package browser

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserApiHandlerPtr = *BrowserApiHandler

type BrowserApiHandler struct {
	browser          Browser
	mcpServerHandler http.Handler
}

func NewBrowserApiHandler(browser Browser, mcpServerHandler func(Browser) http.Handler) *BrowserApiHandler {
	return &BrowserApiHandler{
		browser:          browser,
		mcpServerHandler: mcpServerHandler(browser),
	}
}

func (handler *BrowserApiHandler) NewChiOpenApiRouter(prefix string, opts ...option.OpenAPIOption) chiopenapi.Router {
	opts = stl.NewSlice(
		option.WithTitle("Browser"),
		option.WithDescription("Browser API"),
	).Append(opts...)

	api := NewChiOpenApiRouter(prefix, opts...)

	browserFn := GetBrowserFromRequest(func(r *http.Request) (*BrowserApiHandler, error) {
		return handler, nil
	})

	if prefix == "" || prefix == "/" {
		browserFn.RegisterChiOpenApiRoutes(api)
		return api
	}

	api.Route(prefix, func(r chiopenapi.Router) {
		browserFn.RegisterChiOpenApiRoutes(r)
	})

	return api
}

func (handler *BrowserApiHandler) NewTabHandler(id string) *TabHandler {
	return NewTabHandler(handler.browser, id)
}

func (handler *BrowserApiHandler) HandleGetCdpAddress(_ NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[string] {
	address, err := handler.browser.GetCdpAddress()
	if err != nil {
		return types.NewTypedResponseResultFromError[string](http.StatusInternalServerError, err, "Failed to get browser instance")
	}

	return types.NewTypedResponseResultFromData(address.String())
}

func (handler *BrowserApiHandler) HandleListTabs(_ NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[Tabs] {
	tabs, err := handler.browser.ListTabs()
	if err != nil {
		return types.NewTypedResponseResultFromError[Tabs](http.StatusInternalServerError, err, "Failed to list tabs")
	}

	return types.NewTypedResponseResultFromData(tabs)
}

func (handler *BrowserApiHandler) HandleCreateTab(params *NewTabOptions, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[*Tab] {
	tab, err := handler.browser.NewTab(params)
	if err != nil {
		return types.NewTypedResponseResultFromError[*Tab](http.StatusInternalServerError, err, "Failed to create tab")
	}

	return types.NewTypedResponseResultFromData(tab)
}

type GetBrowserFromRequest func(*http.Request) (*BrowserApiHandler, error)

func (fn GetBrowserFromRequest) RegisterChiOpenApiRoutes(r chiopenapi.Router) {
	r.Route("/Tabs", func(r chiopenapi.Router) {
		RegisterParamsBasedRequestHandler(r, http.MethodGet, "/", BrowserApiHandlerPtr.HandleListTabs, fn).With(
			option.Summary("List"),
			option.Description("List all tabs"),
			option.Tags("Tabs"),
		)

		RegisterParamsBasedRequestHandler(r, http.MethodPost, "/", BrowserApiHandlerPtr.HandleCreateTab, fn).With(
			option.Summary("Create"),
			option.Description("Create a new tab"),
			option.Tags("Tabs"),
		)

		r.Route("/{tabId}", func(r chiopenapi.Router) {
			getTab := GetTabHandlerFromRequest(func(r *http.Request) (*TabHandler, error) {
				handler, err := fn(r)
				if err != nil {
					return nil, err
				}

				return handler.NewTabHandler(chi.URLParam(r, "tabId")), nil
			})

			getTab.RegisterChiOpenApiRoutes(r)
		})
	})

	// todo fix me
	r.Mount("/mcp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, err := fn(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := r.URL.Path
		prefix := path[:strings.Index(path, "/mcp")]

		http.StripPrefix(prefix, handler.mcpServerHandler).ServeHTTP(w, r)
	}))

	RegisterParamsBasedRequestHandler(r, http.MethodGet, "/CdpAddress", BrowserApiHandlerPtr.HandleGetCdpAddress, fn).With(
		option.Summary("Get Cdp Address"),
		option.Description("Get the CDP address of a browser instance"),
		option.Tags("Instances"),
		option.Response(http.StatusOK, new(types.TypedResponseResult[string])),
	)

	r.Mount("/cdp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, err := fn(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		address, err := handler.browser.GetCdpAddress()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		path := r.URL.Path
		index := strings.Index(path, "/cdp")
		prefix := path[:index+len("/cdp")]
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   address.String(),
		})

		http.StripPrefix(prefix, proxy).ServeHTTP(w, r)
	}))
}
