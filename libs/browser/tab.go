package browser

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fasionchan/goutils/types"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type TabHandler struct {
	browser Browser
	id      string
}

func NewTabHandler(browser Browser, id string) *TabHandler {
	return &TabHandler{
		browser: browser,
		id:      id,
	}
}

func (h *TabHandler) GetId() string {
	return h.id
}

func (h *TabHandler) GetBrowser() Browser {
	return h.browser
}

func (h *TabHandler) GetTab() (*Tab, error) {
	return h.browser.GetTab(h.id)
}

func (h *TabHandler) Close() error {
	return h.browser.CloseTab(h.id)
}

func (h *TabHandler) Navigate(url string) error {
	return h.browser.Navigate(h.id, url)
}

func (h *TabHandler) GoBack() error {
	return h.browser.GoBack(h.id)
}

func (h *TabHandler) GoForward() error {
	return h.browser.GoForward(h.id)
}

func (h *TabHandler) Reload() error {
	return h.browser.Reload(h.id)
}

func (h *TabHandler) Click(selector, selectorType, button string, count int) error {
	return h.browser.Click(h.id, selector, selectorType, button, count)
}

func (h *TabHandler) Type(selector, selectorType, text string) error {
	return h.browser.Type(h.id, selector, selectorType, text)
}

func (h *TabHandler) SetInputFiles(selector, selectorType string, files []string) error {
	return h.browser.SetInputFiles(h.id, selector, selectorType, files)
}

func (h *TabHandler) Screenshot(opts *ScreenshotOptions) ([]byte, error) {
	return h.browser.Screenshot(h.id, opts)
}

func (h *TabHandler) GetTexts(selector, selectorType string) (types.Strings, error) {
	return h.browser.GetTexts(h.id, selector, selectorType)
}

func (h *TabHandler) GetHtmls(selector, selectorType string) (types.Strings, error) {
	return h.browser.GetHtmls(h.id, selector, selectorType)
}

func (h *TabHandler) SetCookies(cookies []*http.Cookie) error {
	return h.browser.SetCookies(h.id, cookies)
}

func (h *TabHandler) GetCookies() ([]*http.Cookie, error) {
	return h.browser.GetCookies(h.id)
}

func (h *TabHandler) PrintToPdf() (io.ReadCloser, error) {
	return h.browser.PrintToPdf(h.id)
}

func (h *TabHandler) StartScreencast(opts *ScreencastOptions) (*ScreencastStream, error) {
	return h.browser.StartScreencast(h.id, opts)
}

type GetTabHandlerFromRequest func (*http.Request) (*TabHandler, error)

func (fn GetTabHandlerFromRequest) RegisterChiOpenApiRoutes(r chiopenapi.Router) {
	type NavigateOptions struct {
		Url string `json:"url"`
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		tab, err := tabHandler.GetTab()
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get tab").WriteHttpResponse(w)
			return
		}

		types.NewTypedResponseResultFromData(tab).WriteHttpResponse(w)
	}).With(
		option.Description("Get the current tab"),
		option.Tags("Tabs"),
		option.Response(http.StatusOK, new(Tab)),
	)

	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}


		if err := tabHandler.Close(); err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to close page").WriteHttpResponse(w)
			return
		}

		types.NewResponseResultFromData(nil).WriteHttpResponse(w)
	}).With(
		option.Description("Close the current page"),
		option.Tags("Tabs"),
		option.Response(http.StatusOK, nil),
	)

	r.Post("/_navigate", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		var options NavigateOptions
		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to decode request body").WriteHttpResponse(w)
			return
		}

		if err := tabHandler.Navigate(options.Url); err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to navigate").WriteHttpResponse(w)
			return
		}

		types.NewResponseResultFromData(nil).WriteHttpResponse(w)
	}).With(
		option.Description("Navigate to a URL"),
		option.Tags("Navigation"),
		option.Request(new(NavigateOptions)),
	)

	r.Post("/_reload", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		if err := tabHandler.Reload(); err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to reload").WriteHttpResponse(w)
			return
		}

		types.NewResponseResultFromData(nil).WriteHttpResponse(w)
	}).With(
		option.Description("Reload the current page"),
		option.Tags("Navigation"),
	)

	r.Post("/_goBack", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}


		if err := tabHandler.GoBack(); err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to go back").WriteHttpResponse(w)
			return
		}

		types.NewResponseResultFromData(nil).WriteHttpResponse(w)
	}).With(
		option.Description("Go back to the previous page"),
		option.Tags("Navigation"),
	)

	r.Post("/_goForward", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}
		
		if err := tabHandler.GoForward(); err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to go forward").WriteHttpResponse(w)
			return
		}

		types.NewResponseResultFromData(nil).WriteHttpResponse(w)
	}).With(
		option.Description("Go forward to the next page"),
		option.Tags("Navigation"),
	)

	r.Get("/Screenshot", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		opts, err := NewScreenshotOptionsFromUrlValues(r.URL.Query())
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to parse screenshot options").WriteHttpResponse(w)
			return
		}

		screenshot, err := tabHandler.Screenshot(opts)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to take screenshot").WriteHttpResponse(w)
			return
		}

		w.Header().Set("Content-Type", opts.MimeType())
		w.WriteHeader(http.StatusOK)
		w.Write(screenshot)
	}).With(
		option.Description("Take a screenshot of the current page"),
		option.Tags("Screenshots"),
		option.Response(http.StatusOK, new(bytes.Buffer)),
	)

	r.Route("/Cookies", func(r chiopenapi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			tabHandler, err := fn(r)
			if err != nil {
				types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
				return
			}

			cookies, err := tabHandler.GetCookies()
			if err != nil {
				types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get cookies").WriteHttpResponse(w)
				return
			}

			types.NewTypedResponseResultFromData(cookies).WriteHttpResponse(w)
		}).With(
			option.Description("Get cookies"),
			option.Tags("Cookies"),
			option.Response(http.StatusOK, new([]http.Cookie)),
		)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			tabHandler, err := fn(r)
			if err != nil {
				types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
				return
			}

			var cookie http.Cookie
			if err := json.NewDecoder(r.Body).Decode(&cookie); err != nil {
				types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to decode request body").WriteHttpResponse(w)
				return
			}

			if err := tabHandler.SetCookies([]*http.Cookie{&cookie}); err != nil {
				types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to set cookies").WriteHttpResponse(w)
				return
			}

			types.NewResponseResultFromData(nil).WriteHttpResponse(w)
		}).With(
			option.Description("Set cookies"),
			option.Tags("Cookies"),
			option.Request(new(http.Cookie)),
		)
	})

	r.Get("/Texts", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		request, err := ParseRequest[TabGetTextsRequest](r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to parse request").WriteHttpResponse(w)
			return
		}

		texts, err := tabHandler.GetTexts(request.Target, request.TargetType)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get texts").WriteHttpResponse(w)
			return
		}

		types.NewTypedResponseResultFromData(texts).WriteHttpResponse(w)
	}).With(
		option.Description("Get texts"),
		option.Tags("Texts"),
		option.Request(new(TabGetTextsRequest)),
		option.Response(http.StatusOK, new(types.Strings)),
	)

	r.Get("/Htmls", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		request, err := ParseRequest[TabGetHtmlsRequest](r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to parse request").WriteHttpResponse(w)
			return
		}

		htmls, err := tabHandler.GetHtmls(request.Target, request.TargetType)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get htmls").WriteHttpResponse(w)
			return
		}

		types.NewTypedResponseResultFromData(htmls).WriteHttpResponse(w)
	}).With(
		option.Description("Get htmls"),
		option.Tags("Htmls"),
		option.Request(new(TabGetHtmlsRequest)),
		option.Response(http.StatusOK, new(types.Strings)),
	)
}

type TabGetHtmlsRequest = TabGetTextsRequest

type TabGetTextsRequest struct {
	TabGetTextsQuery `json:"-" query:",inline"`
}

type TabGetTextsQuery struct {
	Target string `query:"target"`
	TargetType string `query:"targetType"`
}