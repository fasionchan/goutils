package browser

import (
	"context"
	"net/http"
	"sync"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserPool struct {
	builder BrowserBuilder
	browsers *stl.SyncMap[string, Browser]
}

func NewBrowserPool(builder BrowserBuilder) *BrowserPool {
	return &BrowserPool{
		builder: builder,
		browsers: stl.NewSyncMap[string, Browser](),
	}
}

func (p *BrowserPool) EnsureBrowser(ctx context.Context, id string) (Browser, error) {
	browser, _, err := p.browsers.LoadOrCreate(ctx, id, func(ctx context.Context) (Browser, error) {
		return p.builder.Build()
	})
	return browser, err
}

func (p *BrowserPool) GetChiOpenApiRouter() chiopenapi.Router {
	return sync.OnceValue(p.NewChiOpenApiRouter)()
}

func (p *BrowserPool) NewChiOpenApiRouter() chiopenapi.Router {
	api := chiopenapi.NewRouter(chi.NewRouter())
	p.RegistryChiOpenApiRoutes(api)
	return api
}

func (p *BrowserPool) RegistryChiOpenApiRoutes(r chiopenapi.Router) {
	r.Route("/Instances", func(r chiopenapi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			instances := p.browsers.Keys()
			types.NewTypedResponseResultFromData(instances).WriteHttpResponse(w)
		}).With(
			option.Description("List all browser instances"),
			option.Tags("Instances"),
			option.Response(http.StatusOK, new(types.TypedResponseResult[types.Strings])),
		)

		r.Route("/{instanceId}", func(r chiopenapi.Router) {
			GetBrowserFromRequest(func(r *http.Request) (Browser, error) {
				return p.EnsureBrowser(r.Context(), chi.URLParam(r, "instanceId"))
			}).RegisterChiOpenApiRoutes(r)
		})
	})
}

func (p *BrowserPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.GetChiOpenApiRouter().ServeHTTP(w, r)
}

func (p *BrowserPool) Close() error {
	// todo
	return nil
}