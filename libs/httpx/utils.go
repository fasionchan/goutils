package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fasionchan/goutils/baseutils/netutils"
	"github.com/fasionchan/goutils/libs/logging"
	"github.com/fasionchan/goutils/std/reflectx"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"go.uber.org/zap"
)

const (
	EnvNameOuterProxies = "OUTER_PROXIES"
)

func NewOuterProxyReplacerFromDefaultEnvPro(lookupEnv func(string) (string, bool)) (PrefixesReplacer, error) {
	return NewOuterProxyReplacerFromEnvPro(EnvNameOuterProxies, lookupEnv)
}

func NewOuterProxyReplacerFromDefaultEnv() (PrefixesReplacer, error) {
	return NewOuterProxyReplacerFromEnvPro(EnvNameOuterProxies, os.LookupEnv)
}

func NewOuterProxyReplacerFromEnvPro(envName string, lookupEnv func(string) (string, bool)) (PrefixesReplacer, error) {
	value, ok := lookupEnv(envName)
	if !ok {
		return PrefixesReplacer{}, nil
	}

	return NewPrefixesReplacer(value)
}

func NewOuterProxyReplacerFromEnv(envName string) (PrefixesReplacer, error) {
	return NewOuterProxyReplacerFromEnvPro(envName, os.LookupEnv)
}

type PrefixesReplacer map[string]string

func NewPrefixesReplacer(config string) (PrefixesReplacer, error) {
	repalcer := PrefixesReplacer{}

	config = strings.TrimSpace(config)
	if config == "" {
		return repalcer, nil
	}

	for _, pair := range strings.Split(config, "||||") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.Split(pair, "===>")
		if len(parts) < 2 {
			return nil, fmt.Errorf("bad outer proxy: %s", pair)
		}

		repalcer[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return repalcer, nil
}

func (replacer PrefixesReplacer) Replace(s string) string {
	if replacer == nil {
		return s
	}

	for prefix, substitute := range replacer {
		if strings.HasPrefix(s, prefix) {
			return substitute + s[len(prefix):]
		}
	}

	return s
}

func (replacer PrefixesReplacer) FilterByKeys(keys types.Strings) PrefixesReplacer {
	result := PrefixesReplacer{}

	if replacer == nil {
		return result
	}

	for key, value := range replacer {
		if keys.Contain(key) {
			result[key] = value
		}
	}

	return result
}

func FormatRange(unit string, start, end int64) string {
	return fmt.Sprintf("%s=%d-%d", unit, start, end)
}

func ParseContentRange(value string) (unit string, start, end int64, total int64, err error) {
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return "", 0, 0, 0, fmt.Errorf("invalid content range: %s", value)
	}

	unit = parts[0]

	parts = strings.SplitN(parts[1], "/", 2)
	if len(parts) != 2 {
		return "", 0, 0, 0, fmt.Errorf("invalid content range: %s", value)
	}

	totalPart := strings.TrimSpace(parts[1])
	if totalPart == "*" {
		total = -1
	} else {
		total, err = strconv.ParseInt(totalPart, 10, 64)
		if err != nil {
			return "", 0, 0, 0, err
		}
	}

	if parts[0] == "*" {
		start = -1
		end = -1
		return
	}

	parts = strings.SplitN(parts[0], "-", 2)
	if len(parts) != 2 {
		return "", 0, 0, 0, fmt.Errorf("invalid content range: %s", value)
	}

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return "", 0, 0, 0, err
	}

	end, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, 0, 0, err
	}

	return
}

type HttpHandler = http.Handler

func GetHandlerOrNotImplemented(data any) http.Handler {
	handler, ok := data.(http.Handler)
	if !ok {
		return NewNotImplementedHandlerWithDataType(data)
	}
	return handler
}

type NotImplementedHandler = StaticHandler

func NewNotImplementedHandler() StaticHandler {
	return NewNotImplementedHandlerWithMessage("")
}

func NewNotImplementedHandlerWithMessage(message string) StaticHandler {
	if message == "" {
		message = "http.Handler not implemented"
	}

	return StaticHandler{
		StatusCode:  http.StatusNotImplemented,
		ContentType: "text/plain",
		Body:        []byte(message),
	}
}

func NewNotImplementedHandlerWithDataType(data any) http.Handler {
	message := fmt.Sprintf("http.Handler not implemented: %s", reflectx.TypeNameOrString(data))
	return NewNotImplementedHandlerWithMessage(message)
}

type StaticHandler struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func NewStaticHandler(statusCode int, contentType string, body []byte) *StaticHandler {
	return &StaticHandler{
		StatusCode:  statusCode,
		ContentType: contentType,
		Body:        body,
	}
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handler.ContentType)
	w.WriteHeader(handler.StatusCode)
	w.Write(handler.Body)
}

type GetJsonHandlerFunc[Data any] func(ctx context.Context) (Data, error)

func NewGetJsonHandlerFunc[Data any](handler GetJsonHandlerFunc[Data]) GetJsonHandlerFunc[Data] {
	return handler
}

func NewGetJsonHandlerFuncPure[Data any](handler func() Data) GetJsonHandlerFunc[Data] {
	return func(ctx context.Context) (Data, error) {
		return handler(), nil
	}
}

func (handler GetJsonHandlerFunc[Data]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loggerRef := logging.LoggerRefFromContext(r.Context())

	data, err := handler(r.Context())
	if err != nil {
		loggerRef.Error("GetJsonHandlerFunc.ServeHTTP.handler Failed",
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		loggerRef.Error("GetJsonHandlerFunc.ServeHTTP.Encode Failed",
			zap.Error(err),
		)

		return
	}
}

func NewCookiesFromUrlValues(values url.Values) []*http.Cookie {
	return netutils.NewCookiesFromMap(UrlValuesToMap(values, stl.LastOneOrZero[[]string]))
}

func UrlValuesToMap(values url.Values, unique func([]string) string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	result := make(map[string]string, len(values))
	for name, values := range values {
		if len(values) == 0 {
			continue
		}

		result[name] = unique(values)
	}

	return result
}