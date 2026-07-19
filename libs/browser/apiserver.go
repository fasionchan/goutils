package browser

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
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
	api := chiopenapi.NewRouter(router)

	api.Get("/Tabs", func(w http.ResponseWriter, r *http.Request) {
		tabs, err := handler.browser.ListTabs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tabs)
	})

	api.Post("/Tabs", func(w http.ResponseWriter, r *http.Request) {
		var options NewTabOptions
		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tab, err := handler.browser.NewTab(&options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tab)
	}).With(
		option.Request(new(NewTabOptions)),
	)

	api.Route("/Tabs/{id}", func(r chiopenapi.Router) {
		getId := func(r *http.Request) string {
			return chi.URLParam(r, "id")
		}

		NewPageApiHandler(handler.browser, getId, &websocket.Upgrader{}).RegisterChiOpenApiRoutes(r)
	})

	return api
}

func (handler *BrowserApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.GetHttpHandler().ServeHTTP(w, r)
}

type PageApiHandler struct {
	browser  Browser
	getId    func(*http.Request) string
	upgrader *websocket.Upgrader
}

func NewPageApiHandler(browser Browser, getId func(*http.Request) string, upgrader *websocket.Upgrader) *PageApiHandler {
	return &PageApiHandler{
		browser:  browser,
		getId:    getId,
		upgrader: upgrader,
	}
}

// func (handler *PageApiHandler) GetMux() *chi.Mux {
// 	return sync.OnceValue(handler.NewMux)()
// }

// func (handler *PageApiHandler) NewMux() *chi.Mux {
// 	mux := chi.NewRouter()
// 	handler.RegisterRoutes(mux)
// 	return mux
// }

func (handler *PageApiHandler) RegisterChiOpenApiRoutes(r chiopenapi.Router) {
	type NavigateOptions struct {
		Url string `json:"url"`
	}

	r.Post("/_navigate", func(w http.ResponseWriter, r *http.Request) {
		id := handler.getId(r)

		var options NavigateOptions
		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := handler.browser.Navigate(id, options.Url); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).With(
		option.Request(new(NavigateOptions)),
	)

	r.Route("/Cookies", func(r chiopenapi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			id := handler.getId(r)
			cookies, err := handler.browser.GetCookies(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(cookies)
		}).With(
			option.Response(http.StatusOK, new([]http.Cookie)),
		)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			id := handler.getId(r)

			var cookie http.Cookie
			if err := json.NewDecoder(r.Body).Decode(&cookie); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := handler.browser.SetCookies(id, []*http.Cookie{&cookie}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}).With(
			option.Request(new(http.Cookie)),
		)
	})
}

// func (handler *PageApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("PageApiHandler", r.URL.Path)
// 	handler.GetMux().ServeHTTP(w, r)
// }

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
