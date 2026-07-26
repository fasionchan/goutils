package browser

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec-ui/stoplight"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

var (
	queryDecoder = form.NewDecoder()
	queryEncoder = form.NewEncoder()
)

type RequestParams[
	Body any,
	Query any,
	Path any,
] struct {
	Body  Body
	Query Query
	Path  Path
}

func ParseRequest[
	Request any,
](r *http.Request) (result Request, err error) {
	if query := r.URL.Query(); len(query) > 0 {
		if err = queryDecoder.Decode(&result, query); err != nil {
			return
		}
	}

	if _, ok := any(result).(NoRequestBody); ok {
		return
	}

	switch contentType := strings.ToLower(r.Header.Get("Content-Type")); {
	case strings.Contains(contentType, "application/json"):
		if err = json.NewDecoder(r.Body).Decode(&result); err != nil {
			return
		}
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		// todo
		if err = form.NewDecoder().Decode(&result, nil); err != nil {
			return
		}
	}

	return
}

type ParamsBasedRequestHandler[
	TargetFromRequest ~func(*http.Request) (Target, error),
	Result any,
	Target any,
	Params any,
] func(target Target, params Params, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[Result]

func NewParamsBasedRequestHandler[
	TargetFromRequest ~func(*http.Request) (Target, error),
	Result any,
	Target any,
	Params any,
](handler func(target Target, params Params, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[Result]) ParamsBasedRequestHandler[TargetFromRequest, Result, Target, Params] {
	return handler
}

func (handler ParamsBasedRequestHandler[TargetFromRequest, Result, Target, Params]) RegisterToChiOpenApiRouter(r chiopenapi.Router, method, path string, targetFromRequest TargetFromRequest) chiopenapi.Route {
	var params any = new(Params)

	opts := stl.NewSlice(
		option.Request(params),
		option.Response(http.StatusOK, new(types.TypedResponseResult[Result])),
	)

	if _, ok := params.(NoRequestBody); ok {
		opts = opts.Append(option.CustomizeOperation(func(operation *openapi.Operation) {
			operation.RequestBody = nil
		}))
	}

	return r.MethodFunc(method, path, func(w http.ResponseWriter, r *http.Request) {
		target, err := targetFromRequest(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to get target").WriteHttpResponse(w)
			return
		}

		params, err := ParseRequest[Params](r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to parse params").WriteHttpResponse(w)
			return
		}

		handler(target, params, w, r).WriteHttpResponse(w)
	}).With(opts...)
}

func RegisterParamsBasedRequestHandler[
	Handler ~func(target Target, params Params, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[Result],
	TargetFromRequest ~func(*http.Request) (Target, error),
	Result any,
	Params any,
	Target any,
](r chiopenapi.Router, method, path string, handler Handler, targetFromRequest TargetFromRequest) chiopenapi.Route {
	return NewParamsBasedRequestHandler[TargetFromRequest](handler).
		RegisterToChiOpenApiRouter(r, method, path, targetFromRequest)
}

// todo move to httpx
type ResponseWriterBodyCapturer struct {
	http.ResponseWriter
	bytes.Buffer
}

func NewResponseWriterBodyCapturer(w http.ResponseWriter) *ResponseWriterBodyCapturer {
	return &ResponseWriterBodyCapturer{
		ResponseWriter: w,
	}
}

func (w *ResponseWriterBodyCapturer) GetBody() []byte {
	return w.Buffer.Bytes()
}

func (w *ResponseWriterBodyCapturer) Write(p []byte) (n int, err error) {
	return w.Buffer.Write(p)
}

func NewChiOpenApiRouter(basePath string, opts ...option.OpenAPIOption) chiopenapi.Generator {
	if basePath == "" {
		basePath = "/"
	}

	serverUrlMagic := "__server__"
	docsPath := path.Join(basePath, "docs")
	specPath := path.Join(docsPath, "openapi.yaml")

	opts = stl.NewSlice(
		option.WithSpecPath(specPath),
		option.WithDocsPath(docsPath),

		option.WithServer(serverUrlMagic),

		// WithUIOption 会覆盖默认 UI provider，因此这里需显式保留 Stoplight。
		// 服务端路由仍用默认绝对路径 /docs/openapi.yaml（chi 要求以 / 开头）；
		// UI 的 SpecPath 改为相对路径，前端会按当前页面 pathname 拼接，兼容 nginx 前缀改写。
		option.WithUIOption(func(c *config.SpecUI) {
			stoplight.WithUI()(c)
			specui.WithSpecPath("./openapi.yaml")(c)
		}),

		AddDocumentCustomizers(UnwindObjectQueryParametersForDocument),
	).Append(opts...)

	docsApi := chiopenapi.NewRouter(chi.NewRouter(), opts...)
	api := chiopenapi.NewRouter(chi.NewRouter(), opts...)

	api.HandleFunc(docsPath, func(w http.ResponseWriter, r *http.Request) {
		capture := NewResponseWriterBodyCapturer(w)
		docsApi.ServeHTTP(capture, r)

		body := capture.GetBody()

		docsLine := "const docs = document.getElementById('docs');"
		addLine := "docs.tryItCredentialsPolicy = 'same-origin';"

		body = bytes.Replace(body, []byte(docsLine), []byte(docsLine+"    "+addLine), 1)

		w.Write(body)
	})

	api.HandleFunc(specPath, func(w http.ResponseWriter, r *http.Request) {
		spec, err := api.GenerateSchema("yaml")
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to generate OpenAPI schema").WriteHttpResponse(w)
			return
		}

		prefix := r.Header.Get("X-Forwarded-Prefix")
		if prefix == "" {
			prefix = "/"
		}

		old := []byte("- url: " + serverUrlMagic)
		new := []byte("- url: " + path.Clean(prefix))
		spec = bytes.Replace(spec, old, new, 1)

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write(spec)
	})

	return api
}

func AddDocumentCustomizers(customizer func(*openapi.Document)) option.OpenAPIOption {
	return func(c *openapi.Config) {
		c.DocumentCustomizers = append(c.DocumentCustomizers, customizer)
	}
}

func UnwindObjectQueryParametersForDocument(doc *openapi.Document) {
	for _, path := range doc.Paths {
		stl.NewSlice(
			path.Get,
			path.Post,
			path.Put,
			path.Delete,
			path.Options,
			path.Head,
			path.Patch,
			path.Trace,
		).ForEach(UnwindObjectQueryParametersForOperation)
	}
}

func UnwindObjectQueryParametersForOperation(operation *openapi.Operation) {
	if operation == nil {
		return
	}

	operation.Parameters = UnwindObjectQueryParameters(operation.Parameters...)
}

func UnwindObjectQueryParameters(parameters ...*openapi.Parameter) []*openapi.Parameter {
	result := make([]*openapi.Parameter, 0, len(parameters))
	for _, parameter := range parameters {
		if parameter.In != "query" {
			result = append(result, parameter)
			continue
		}
		result = append(result, UnwindObjectParamterInDotFormat(parameter)...)
	}

	return result
}

func UnwindObjectParamterInDotFormat(parameter *openapi.Parameter) []*openapi.Parameter {
	schema := parameter.Schema
	if schema == nil {
		return []*openapi.Parameter{parameter}
	}

	if len(schema.Properties) == 0 {
		return []*openapi.Parameter{parameter}
	}

	result := make([]*openapi.Parameter, 0, len(schema.Properties))
	for name, property := range schema.Properties {
		subParameter := stl.Dup(parameter)
		subParameter.Name = parameter.Name + "." + name
		subParameter.Schema = property
		result = append(result, UnwindObjectParamterInDotFormat(subParameter)...)
	}

	return result
}

type NoRequestBody interface {
	isNoRequestBody()
}

type NoRequestBodyParams struct{}

func (p NoRequestBodyParams) isNoRequestBody() {}

func GenerateJwtPathToken(path, secret string, expireDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"path": path,
		"exp":  time.Now().Add(expireDuration).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func JwtPathTokenAuthorizer(secret string, expireDuration time.Duration) func(*http.Request) error {
	return func(r *http.Request) error {
		path := strings.Trim(r.URL.Path, "/")
		token, err := GenerateJwtPathToken(path, secret, expireDuration)
		if err != nil {
			return err
		}

		r.Header.Set("Authorization", token)

		return nil
	}
}

func JwtPathValidator(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsed, err := jwt.Parse(r.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				types.NewResponseResultFromError(http.StatusUnauthorized, err, "Invalid token").WriteHttpResponse(w)
				return
			}

			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				types.NewResponseResultFromError(http.StatusUnauthorized, errors.New("invalid claims"), "Invalid claims").WriteHttpResponse(w)
				return
			}

			if claims.Valid() != nil {
				types.NewResponseResultFromError(http.StatusUnauthorized, claims.Valid(), "Invalid token").WriteHttpResponse(w)
				return
			}

			path, ok := claims["path"].(string)
			if !ok {
				types.NewResponseResultFromError(http.StatusUnauthorized, errors.New("invalid path type"), "Invalid path type").WriteHttpResponse(w)
				return
			}

			if strings.Trim(path, "/") != strings.Trim(r.URL.Path, "/") {
				types.NewResponseResultFromError(http.StatusUnauthorized, errors.New("invalid path"), "Invalid path").WriteHttpResponse(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func init() {
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}
