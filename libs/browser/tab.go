package browser

import (
	"bytes"
	"io"
	"net/http"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/option"
)

type TabHandlerPtr = *TabHandler

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

func (h *TabHandler) Hover(selector, selectorType string) error {
	return h.browser.Hover(h.id, selector, selectorType)
}

func (h *TabHandler) SelectOption(target, targetType string, options []string, optionType string, selected bool) error {
	return h.browser.SelectOption(h.id, target, targetType, options, optionType, selected)
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

func (h *TabHandler) NewRemoteController() *RemoteController {
	return NewRemoteController(h.browser, h.id, DefaultWebSocketUpgrader)
}

func (h *TabHandler) HandleGetTab(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[*Tab] {
	tab, err := h.GetTab()
	if err != nil {
		return types.NewTypedResponseResultFromError[*Tab](http.StatusBadRequest, err, "Failed to get tab")
	}

	return types.NewTypedResponseResultFromData(tab)
}

func (h *TabHandler) HandleDeleteTab(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.ResponseResult {
	if err := h.Close(); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to close tab")
	}

	return types.NewResponseResultFromData(nil)
}

type NavigateOptions struct {
	Url string `json:"url"`
}

func (h *TabHandler) HandleNavigate(params *NavigateOptions, w http.ResponseWriter, r *http.Request) *types.ResponseResult {
	if err := h.Navigate(params.Url); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to navigate")
	}

	return types.NewResponseResultFromData(nil)
}

func (h *TabHandler) HandleReload(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.ResponseResult {
	if err := h.Reload(); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to reload")
	}

	return types.NewResponseResultFromData(nil)
}

func (h *TabHandler) HandleGoBack(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.ResponseResult {
	if err := h.GoBack(); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to go back")
	}

	return types.NewResponseResultFromData(nil)
}

func (h *TabHandler) HandleGoForward(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.ResponseResult {
	if err := h.GoForward(); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to go forward")
	}

	return types.NewResponseResultFromData(nil)
}

func (h *TabHandler) HandleListCookies(params NoRequestBodyParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[[]*http.Cookie] {
	cookies, err := h.GetCookies()
	if err != nil {
		return types.NewTypedResponseResultFromError[[]*http.Cookie](http.StatusInternalServerError, err, "Failed to get cookies")
	}

	return types.NewTypedResponseResultFromData(cookies)
}

func (h *TabHandler) HandleSetCookie(params *http.Cookie, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[*http.Cookie] {
	if err := h.SetCookies([]*http.Cookie{params}); err != nil {
		return types.NewTypedResponseResultFromError[*http.Cookie](http.StatusInternalServerError, err, "Failed to set cookie")
	}

	cookies, err := h.GetCookies()
	if err != nil {
		return types.NewTypedResponseResultFromError[*http.Cookie](http.StatusInternalServerError, err, "Failed to get cookies")
	}

	cookie, _ := stl.FindFirstByKey(cookies, func(cookie *http.Cookie) string {
		return cookie.Name
	}, params.Name)

	return types.NewTypedResponseResultFromData(cookie)
}

type SnapshotOptions struct {
	Type *string `json:"type"`
}

func (opts *SnapshotOptions) GetType() string {
	if opts == nil {
		return ""
	}

	if ptr := opts.Type; ptr != nil {
		return *ptr
	}

	return ""
}

func (h *TabHandler) HandleSnapshot(params *SnapshotOptions, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[string] {
	snapshot, err := h.browser.Snapshot(h.id, params.GetType())
	if err != nil {
		return types.NewTypedResponseResultFromError[string](http.StatusInternalServerError, err, "Failed to take snapshot")
	}

	return types.NewTypedResponseResultFromData(snapshot)
}

type TabClickRequestParams struct {
	Target     string `json:"target"`
	TargetType string `json:"targetType"`
	Button     string `json:"button"`
	Count      int    `json:"count"`
}

func (h *TabHandler) HandleClick(params *TabClickRequestParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[any] {
	if err := h.Click(params.Target, params.TargetType, params.Button, params.Count); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to click")
	}

	return types.NewResponseResultFromData(nil)
}

type TabTypeRequestParams struct {
	Target     string `json:"target"`
	TargetType string `json:"targetType"`
	Text       string `json:"text"`
}

func (h *TabHandler) HandleType(params *TabTypeRequestParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[any] {
	if err := h.Type(params.Target, params.TargetType, params.Text); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to type")
	}

	return types.NewResponseResultFromData(nil)
}

type TabHoverRequestParams struct {
	Target     string `json:"target"`
	TargetType string `json:"targetType"`
}

func (h *TabHandler) HandleHover(params *TabHoverRequestParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[any] {
	if err := h.Hover(params.Target, params.TargetType); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to hover")
	}

	return types.NewResponseResultFromData(nil)
}

type TabSelectOptionRequestParams struct {
	Target     string   `json:"target"`
	TargetType string   `json:"targetType"`
	Options    []string `json:"options"`
	OptionType string   `json:"optionType"`
	Selected   bool     `json:"selected"`
}

func (h *TabHandler) HandleSelectOption(params *TabSelectOptionRequestParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[any] {
	if err := h.SelectOption(params.Target, params.TargetType, params.Options, params.OptionType, params.Selected); err != nil {
		return types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to select option")
	}

	return types.NewResponseResultFromData(nil)
}

type TabGetHtmlsParams struct {
	DomElementLocatorQuery `json:"-" query:",inline"`
}

func (h *TabHandler) HandleGetHtmls(params *TabGetHtmlsParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[types.Strings] {
	htmls, err := h.GetHtmls(params.Target, params.TargetType)
	if err != nil {
		return types.NewTypedResponseResultFromError[types.Strings](http.StatusInternalServerError, err, "Failed to get htmls")
	}

	return types.NewTypedResponseResultFromData(htmls)
}

type TabGetTextsParamsPtr = *TabGetTextsParams

type TabGetTextsParams struct {
	DomElementLocatorQuery `json:"-" query:",inline"`
}

func (h *TabHandler) HandleGetTexts(params *TabGetTextsParams, w http.ResponseWriter, r *http.Request) *types.TypedResponseResult[types.Strings] {
	texts, err := h.GetTexts(params.Target, params.TargetType)
	if err != nil {
		return types.NewTypedResponseResultFromError[types.Strings](http.StatusInternalServerError, err, "Failed to get texts")
	}

	return types.NewTypedResponseResultFromData(texts)
}

type DomElementLocatorQuery struct {
	Target     string `query:"target"`
	TargetType string `query:"targetType"`
}

type GetTabHandlerFromRequest func(*http.Request) (*TabHandler, error)

func (fn GetTabHandlerFromRequest) RegisterChiOpenApiRoutes(r chiopenapi.Router) {

	RegisterParamsBasedRequestHandler(r, http.MethodGet, "/", TabHandlerPtr.HandleGetTab, fn).With(
		option.Summary("Get"),
		option.Description("Get the current tab"),
		option.Tags("Tabs"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodDelete, "/", TabHandlerPtr.HandleDeleteTab, fn).With(
		option.Summary("Close"),
		option.Description("Close the current page"),
		option.Tags("Tabs"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_navigate", TabHandlerPtr.HandleNavigate, fn).With(
		option.Summary("Navigate"),
		option.Description("Navigate to a URL"),
		option.Tags("Navigation"),
		option.Request(new(NavigateOptions)),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_reload", TabHandlerPtr.HandleReload, fn).With(
		option.Summary("Reload"),
		option.Description("Reload the current page"),
		option.Tags("Navigation"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_goBack", TabHandlerPtr.HandleGoBack, fn).With(
		option.Summary("Go back"),
		option.Description("Go back to the previous page"),
		option.Tags("Navigation"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_goForward", TabHandlerPtr.HandleGoForward, fn).With(
		option.Summary("Go forward"),
		option.Description("Go forward to the next page"),
		option.Tags("Navigation"),
	)

	r.Get("/Screenshot", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		query, err := ParseRequest[ScreenOptionsQuery](r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusBadRequest, err, "Failed to parse screenshot options").WriteHttpResponse(w)
			return
		}

		opts := query.ToScreenshotOptions()

		screenshot, err := tabHandler.Screenshot(opts)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to take screenshot").WriteHttpResponse(w)
			return
		}

		w.Header().Set("Content-Type", opts.MimeType())
		w.WriteHeader(http.StatusOK)
		w.Write(screenshot)
	}).With(
		option.Summary("Screenshot"),
		option.Description("Take a screenshot of the current page"),
		option.Tags("Observations"),
		option.Request(new(ScreenOptionsQuery)),
		option.Response(http.StatusOK, new(bytes.Buffer)),
	)

	r.Route("/Cookies", func(r chiopenapi.Router) {
		RegisterParamsBasedRequestHandler(r, http.MethodGet, "/", TabHandlerPtr.HandleListCookies, fn).With(
			option.Summary("List"),
			option.Description("List cookies"),
			option.Tags("Cookies"),
		)

		RegisterParamsBasedRequestHandler(r, http.MethodPost, "/", TabHandlerPtr.HandleSetCookie, fn).With(
			option.Summary("Set"),
			option.Description("Set cookie"),
			option.Tags("Cookies"),
		)
	})

	RegisterParamsBasedRequestHandler(r, http.MethodGet, "/Snapshot", TabHandlerPtr.HandleSnapshot, fn).With(
		option.Summary("Snapshot"),
		option.Description("Take a snapshot of the current page"),
		option.Tags("Observations"),
		option.Request(new(SnapshotOptions)),
		option.Response(http.StatusOK, new(string)),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_click", TabHandlerPtr.HandleClick, fn).With(
		option.Summary("Click"),
		option.Description("Click on a element"),
		option.Tags("Actions"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_type", TabHandlerPtr.HandleType, fn).With(
		option.Summary("Type"),
		option.Description("Type text into a element"),
		option.Tags("Actions"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_hover", TabHandlerPtr.HandleHover, fn).With(
		option.Summary("Hover"),
		option.Description("Hover over a element"),
		option.Tags("Actions"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodPost, "/_selectOption", TabHandlerPtr.HandleSelectOption, fn).With(
		option.Summary("Select option"),
		option.Description("Select an option from a dropdown"),
		option.Tags("Actions"),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodGet, "/Texts", TabHandlerPtr.HandleGetTexts, fn).With(
		option.Summary("Texts"),
		option.Description("Get texts"),
		option.Tags("Observations"),
		option.Response(http.StatusOK, new(types.Strings)),
	)

	RegisterParamsBasedRequestHandler(r, http.MethodGet, "/Htmls", TabHandlerPtr.HandleGetHtmls, fn).With(
		option.Summary("Htmls"),
		option.Description("Get htmls"),
		option.Tags("Observations"),
		option.Response(http.StatusOK, new(types.Strings)),
	)

	fn.RegisterRemoteChiOpenApiRoute(r)
}

func (fn GetTabHandlerFromRequest) RegisterRemoteChiOpenApiRoute(r chiopenapi.Router) {
	r.Get("/Remote", func(w http.ResponseWriter, r *http.Request) {
		tabHandler, err := fn(r)
		if err != nil {
			types.NewResponseResultFromError(http.StatusInternalServerError, err, "Failed to get page").WriteHttpResponse(w)
			return
		}

		tabHandler.NewRemoteController().ServeHTTP(w, r)
	}).With(
		option.Summary("Remote control"),
		option.Description("WebSocket remote control: binary screencast frames and JSON mouse/keyboard events"),
		option.Tags("Remote"),
		option.Request(new(ScreencastOptionsQuery)),
	)
}
