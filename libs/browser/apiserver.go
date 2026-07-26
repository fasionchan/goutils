package browser

import (
	"net/http"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserApiHandlerPtr = *BrowserApiHandler

type BrowserApiHandler struct {
	browser Browser
}

func NewBrowserApiHandler(browser Browser) *BrowserApiHandler {
	return &BrowserApiHandler{
		browser: browser,
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
}
