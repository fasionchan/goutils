package browser

import (
	"context"
	"net/http"
	"sync"

	"github.com/fasionchan/goutils/std/netx"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec-ui/stoplight"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserPool struct {
	launcher BrowserLauncher
	opts *BrowserLaunchOptions
	browsers *stl.SyncMap[string, Browser]
}

func NewBrowserPool(opts *BrowserLaunchOptions, launcher BrowserLauncher) *BrowserPool {
	return &BrowserPool{
		opts: opts,
		launcher: launcher,
		browsers: stl.NewSyncMap[string, Browser](),
	}
}

func NewBrowserPoolFromLaunchFunc(opts *BrowserLaunchOptions, launchFunc func(ctx context.Context, opts *BrowserLaunchOptions) (Browser, error)) *BrowserPool {
	return NewBrowserPool(opts, BrowserLaunchFunc(launchFunc))
}

func NewBrowserPoolFromTypedLaunchFunc[
	BrowserT Browser,
](opts *BrowserLaunchOptions, launchFunc func(ctx context.Context, opts *BrowserLaunchOptions) (BrowserT, error)) *BrowserPool {
	return NewBrowserPool(opts, BrowserLaunchFunc(func(ctx context.Context, opts *BrowserLaunchOptions) (Browser, error) {
		return launchFunc(ctx, opts)
	}))
}

func (p *BrowserPool) EnsureBrowser(ctx context.Context, id string) (Browser, error) {
	browser, _, err := p.browsers.LoadOrCreate(ctx, id, func(ctx context.Context) (Browser, error) {
		opts := p.opts.Dup().WithAddr(netx.RandomLocalTcpAddr())
		return p.launcher.Launch(ctx, opts)
	})
	return browser, err
}

func (p *BrowserPool) GetChiOpenApiRouter() chiopenapi.Router {
	return sync.OnceValue(p.NewChiOpenApiRouter)()
}

func (p *BrowserPool) NewChiOpenApiRouter() chiopenapi.Router {
	api := chiopenapi.NewRouter(chi.NewRouter(),
		// WithUIOption 会覆盖默认 UI provider，因此这里需显式保留 Stoplight。
		// 服务端路由仍用默认绝对路径 /docs/openapi.yaml（chi 要求以 / 开头）；
		// UI 的 SpecPath 改为相对路径，前端会按当前页面 pathname 拼接，兼容 nginx 前缀改写。
		option.WithUIOption(func(c *config.SpecUI) {
			stoplight.WithUI()(c)
			specui.WithSpecPath("./openapi.yaml")(c)
		}),
	)
	p.RegistryChiOpenApiRoutes(api)
	return api
}

func (p *BrowserPool) RegistryChiOpenApiRoutes(r chiopenapi.Router) {
	r.Route("/Instances", func(r chiopenapi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			instances := p.browsers.Keys()
			types.NewTypedResponseResultFromData(instances).WriteHttpResponse(w)
		}).With(
			option.Summary("List"),
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