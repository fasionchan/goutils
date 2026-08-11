package httpx

import (
	"net/http"

	"github.com/fasionchan/goutils/libs/logging"
	"github.com/fasionchan/goutils/stl"
)

type HttpMiddleware = Middleware

type Middleware interface {
	Apply(next http.Handler) http.Handler
}

type HttpMiddlewares = Middlewares

type Middlewares []Middleware

func NewMiddlewares(middlewares ...Middleware) Middlewares {
	return middlewares
}

func NewMiddlewaresFromFuncs(funcs ...MiddlewareFunc) Middlewares {
	return stl.Map(funcs, func(fn MiddlewareFunc) Middleware {
		return fn
	})
}

func (middlewares Middlewares) Append(more ...Middleware) Middlewares {
	return append(middlewares, more...)
}

func (middlewares Middlewares) Concat(more ...Middleware) Middlewares {
	return append(middlewares, more...)
}

func (middlewares Middlewares) Apply(next http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i].Apply(next)
	}
	return next
}

type HttpMiddlewareFunc = MiddlewareFunc

type MiddlewareFunc func(next http.Handler) http.Handler

func (middleware MiddlewareFunc) Apply(next http.Handler) http.Handler {
	return middleware(next)
}

type LoggerRefMiddleware string

func NewLoggerRefMiddleware(name string) LoggerRefMiddleware {
	return LoggerRefMiddleware(name)
}

func (middleware LoggerRefMiddleware) Apply(next http.Handler) http.Handler {
	name := string(middleware)
	if name == "" {
		name = "LoggerRefMiddleware"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ctx := logging.LoggerRefFromContextPro(r.Context(), true, true, name, logging.GetLogger().Named(name))
		if ctx != r.Context() {
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
