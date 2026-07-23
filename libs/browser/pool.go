package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fasionchan/goutils/std/netx"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserPool struct {
	launcher BrowserLauncher
	opts     *BrowserLaunchOptions
	browsers *stl.SyncMap[string, Browser]
}

func NewBrowserPool(opts *BrowserLaunchOptions, launcher BrowserLauncher) *BrowserPool {
	return &BrowserPool{
		opts:     opts,
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

func (p *BrowserPool) DeleteBrowser(ctx context.Context, id string) (Browser, error) {
	browser, loaded := p.browsers.Delete(id)
	if !loaded {
		return nil, fmt.Errorf("browser not found")
	}
	err := browser.Close()
	return browser, err
}

func (p *BrowserPool) NewChiOpenApiRouter(prefix string) chiopenapi.Router {
	api := NewChiOpenApiRouter(prefix,
		option.WithTitle("Browser Pool"),
		option.WithDescription("Browser Pool API"),
	)

	if prefix == "" || prefix == "/" {
		p.RegistryChiOpenApiRoutes(api)
		return api
	}

	api.Route(prefix, func(r chiopenapi.Router) {
		p.RegistryChiOpenApiRoutes(r)
	})

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

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			if err := p.Close(); err != nil {
				types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to delete browser instances").WriteHttpResponse(w)
				return
			}
			types.NewResponseResultFromData(nil).WriteHttpResponse(w)
		}).With(
			option.Summary("Delete All"),
			option.Description("Delete all browser instances"),
			option.Tags("Instances"),
			option.Response(http.StatusOK, new(types.TypedResponseResult[types.Strings])),
		)

		r.Route("/{instanceId}", func(r chiopenapi.Router) {
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				instanceId := chi.URLParam(r, "instanceId")
				if _, err := p.DeleteBrowser(r.Context(), instanceId); err != nil {
					types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to delete browser instance").WriteHttpResponse(w)
					return
				}
				types.NewTypedResponseResultFromData(instanceId).WriteHttpResponse(w)
			}).With(
				option.Summary("Delete"),
				option.Description("Delete a browser instance"),
				option.Tags("Instances"),
				option.Response(http.StatusOK, new(types.TypedResponseResult[string])),
			)

			GetBrowserFromRequest(func(r *http.Request) (Browser, error) {
				return p.EnsureBrowser(r.Context(), chi.URLParam(r, "instanceId"))
			}).RegisterChiOpenApiRoutes(r)
		})
	})
}

func (p *BrowserPool) Close() error {
	for _, key := range p.browsers.Keys() {
		if _, err := p.DeleteBrowser(context.Background(), key); err != nil {
			return err
		}
	}
	return nil
}
