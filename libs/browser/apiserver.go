package browser

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/fasionchan/goutils/types"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec-ui/stoplight"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type BrowserApiHandler struct {
	browser Browser
}

func NewBrowserApiHandler(browser Browser) *BrowserApiHandler {
	return &BrowserApiHandler{
		browser: browser,
	}
}

func (handler *BrowserApiHandler) GetHttpHandler() http.Handler {
	return sync.OnceValue(handler.NewHttpHandler)()
}

func (handler *BrowserApiHandler) NewHttpHandler() http.Handler {
	router := chi.NewRouter()

	api := chiopenapi.NewRouter(router,
		option.WithServer("../", option.ServerDescription("Relative to docs (reverse-proxy friendly)")),
		option.WithUIOption(func(c *config.SpecUI) {
			stoplight.WithUI()(c)
			specui.WithSpecPath("./openapi.yaml")(c)
		}),
	)

	GetBrowserFromRequest(func(r *http.Request) (Browser, error) {
		return handler.browser, nil
	}).RegisterChiOpenApiRoutes(api)

	return router
}

func (handler *BrowserApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.GetHttpHandler().ServeHTTP(w, r)
}

type GetBrowserFromRequest func (*http.Request) (Browser, error)

func (fn GetBrowserFromRequest) RegisterChiOpenApiRoutes(r chiopenapi.Router) {
	r.Route("/Tabs", func(r chiopenapi.Router) {
		r.Get("/", fn.listTabs).With(
			option.Summary("List"),
			option.Description("List all tabs"),
			option.Tags("Tabs"),
			option.Response(http.StatusOK, new(Tabs)),
		)

		r.Post("/", fn.createTab).With(
			option.Summary("Create"),
			option.Description("Create a new tab"),
			option.Tags("Tabs"),
			option.Request(new(NewTabOptions)),
			option.Response(http.StatusOK, new(Tab)),
		)

		r.Route("/{tabId}", func(r chiopenapi.Router) {
			getTab := func(r *http.Request) (*TabHandler, error) {
				browser, err := fn(r)
				if err != nil {
					return nil, err
				}

				return NewTabHandler(browser, chi.URLParam(r, "tabId")), nil
			}

			GetTabHandlerFromRequest(getTab).RegisterChiOpenApiRoutes(r)
		})
	})
}

func (fn GetBrowserFromRequest) listTabs(w http.ResponseWriter, r *http.Request) {
	browser, err := fn(r)
	if err != nil {
		types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get browser").WriteHttpResponse(w)
		return
	}

	tabs, err := browser.ListTabs()
	if err != nil {
		types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to list tabs").WriteHttpResponse(w)
		return
	}

	types.NewTypedResponseResultFromData(tabs).WriteHttpResponse(w)
}

func (fn GetBrowserFromRequest) createTab(w http.ResponseWriter, r *http.Request) {
	browser, err := fn(r)
	if err != nil {
		types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get browser").WriteHttpResponse(w)
		return
	}

	var options NewTabOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to decode request body").WriteHttpResponse(w)
		return
	}

	tab, err := browser.NewTab(&options)
	if err != nil {
		types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to create new tab").WriteHttpResponse(w)
		return
	}

	types.NewTypedResponseResultFromData(tab).WriteHttpResponse(w)
}

type RemoteController struct {
	browser  Browser
	id       string
	upgrader *websocket.Upgrader
}

func NewRemoteController(browser Browser, id string, upgrader *websocket.Upgrader) *RemoteController {
	return &RemoteController{
		browser:  browser,
		id:       id,
		upgrader: upgrader,
	}
}

func (controller *RemoteController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opts, err := NewScreencastOptionsFromUrlValues(r.URL.Query())
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := controller.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	frames, err := controller.browser.StartScreencast(controller.id, opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer frames.Close()

	go func() {
		for frame := range frames.BytesChan {
			if err := conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		if err := conn.ReadJSON(&json.RawMessage{}); err != nil {
			log.Println(err)
			return
		}
	}
}