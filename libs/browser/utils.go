package browser

import (
	"encoding/json"
	"net/http"

	"github.com/fasionchan/goutils/types"
	"github.com/go-playground/form/v4"
	"github.com/oaswrap/spec/adapter/chiopenapi"
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

	if err = queryDecoder.Decode(&result, r.URL.Query()); err != nil {
		return
	}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err = json.NewDecoder(r.Body).Decode(&result); err != nil {
			return
		}
	case "application/x-www-form-urlencoded":
		if err = form.NewDecoder().Decode(&result, r.URL.Query()); err != nil {
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
	}).With(
		option.Request(new(Params)),
		option.Response(http.StatusOK, new(types.TypedResponseResult[Result])),
	)
}

func RegisterParamsBasedRequestHandler[
	Handler ~func(target Target, params Params, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[Result],
	TargetFromRequest ~func(*http.Request) (Target, error),
	Result any,
	Params any,
	Target any,
](r chiopenapi.Router, method, path string, handler Handler, targetFromRequest TargetFromRequest) chiopenapi.Route {
	return NewParamsBasedRequestHandler[TargetFromRequest, Result, Target, Params](handler).
		RegisterToChiOpenApiRouter(r, method, path, targetFromRequest)
}

func init() {
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}
