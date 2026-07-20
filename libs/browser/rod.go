package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gorilla/websocket"
)

type RodBrowser rod.Browser

func NewRodBrowser() *RodBrowser {
	return (*RodBrowser)(rod.New())
}

func ConnectRodBrowser() (*RodBrowser, error) {
	browser := NewRodBrowser()
	if err := browser.Native().Connect(); err != nil {
		return nil, err
	}
	return browser, nil
}

func (b *RodBrowser) Native() *rod.Browser {
	return (*rod.Browser)(b)
}

func (b *RodBrowser) Close() error {
	return b.Native().Close()
}

func (b *RodBrowser) NewTab(options *NewTabOptions) (*Tab, error) {
	if options == nil {
		options = new(NewTabOptions)
	}

	page, err := b.Native().Page(proto.TargetCreateTarget{
		URL:    options.Url,
		Width:  options.Width,
		Height: options.Height,
	})
	if err != nil {
		return nil, err
	}

	return RodPageToTab(page), nil
}

func (b *RodBrowser) GetTab(id string) (*Tab, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}
	return RodPageToTab(page), nil
}

func (b *RodBrowser) ListTabs() (Tabs, error) {
	pages, err := b.Native().Pages()
	if err != nil {
		return nil, err
	}
	return RodPagesToTabs(pages), nil
}

func (b *RodBrowser) SwitchToTab(id string) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}

	_, err = page.Activate()
	return err
}

func (b *RodBrowser) CloseTab(id string) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	return page.Close()
}

func (b *RodBrowser) Navigate(id string, url string) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	return page.Navigate(url)
}

func (b *RodBrowser) GoBack(id string) error {
	b.Native().Incognito()
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	return page.NavigateBack()
}

func (b *RodBrowser) GoForward(id string) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	return page.NavigateForward()
}

func (b *RodBrowser) Reload(id string) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	return page.Reload()
}

// func (b *RodBrowser) Wait(id string, timeout time.Duration) error {
// 	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
// 	if err != nil {
// 		return err
// 	}
// 	return page.Wait(timeout)
// }

func (b *RodBrowser) Click(id string, selector, selectorType, button string, count int) error {
	element, err := b.getExactElement(id, selector, selectorType)
	if err != nil {
		return err
	}

	return element.Click(RodMouseButtonFromStd(button), count)
}

func (b *RodBrowser) Type(id string, selector, selectorType, text string) error {
	element, err := b.getExactElement(id, selector, selectorType)
	if err != nil {
		return err
	}

	return element.Input(text)
}

func (b *RodBrowser) SetInputFiles(id string, selector, selectorType string, files []string) error {
	element, err := b.getExactElement(id, selector, selectorType)
	if err != nil {
		return err
	}

	return element.SetFiles(files)
}

func (b *RodBrowser) Screenshot(id string, opts *ScreenshotOptions) ([]byte, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}

	return page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormat(opts.GetFormat()),
		Quality: opts.Quality,
	})
}

// func (b *RodBrowser) Snapshot(id, snapshotType string) (string, error) {
// 	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
// 	if err != nil {
// 		return "", err
// 	}
// 	return page.Snapshot(snapshotType)
// }

func (b *RodBrowser) GetTexts(id, selector, selectorType string) (types.Strings, error) {
	return b.getStrings(id, selector, selectorType, (*rod.Element).Text)
}

func (b *RodBrowser) GetHtmls(id, selector, selectorType string) (types.Strings, error) {
	return b.getStrings(id, selector, selectorType, (*rod.Element).HTML)
}

func (b *RodBrowser) getStrings(id, selector, selectorType string, getter func(*rod.Element) (string, error)) (types.Strings, error) {
	elements, err := b.getElements(id, selector, selectorType)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("no elements found for selector: %s", selector)
	}

	return GetStringsFromRodElements(elements, getter)
}

func (b *RodBrowser) getElements(id, selector, selectorType string) (rod.Elements, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}

	var fn func(string) (rod.Elements, error)
	switch selectorType {
	case SelectorTypeCss:
		fn = page.Elements
	case SelectorTypeXPath:
		fn = page.ElementsX
	default:
		return nil, fmt.Errorf("invalid selector type: %s", selectorType)
	}

	return fn(selector)
}

func (b *RodBrowser) getExactElement(id, selector, selectorType string) (*rod.Element, error) {
	elements, err := b.getElements(id, selector, selectorType)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("no elements found for selector: %s", selector)
	} else if len(elements) > 1 {
		return nil, fmt.Errorf("multiple elements found for selector: %s", selector)
	}

	return elements.First(), nil
}

func (b *RodBrowser) GetCookies(id string) ([]*http.Cookie, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}

	cookies, err := page.Cookies([]string{})
	if err != nil {
		return nil, err
	}

	return NetworkCookiesToStd(cookies), nil
}

// SameSite=None requires Secure=true
func (b *RodBrowser) SetCookies(id string, cookies []*http.Cookie) error {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return err
	}
	json.NewEncoder(os.Stdout).Encode(NetworkCookieParamsFromStd(cookies))
	fmt.Println("--------------------------------")
	return page.SetCookies(NetworkCookieParamsFromStd(cookies))
}

func (b *RodBrowser) PrintToPdf(id string) (io.ReadCloser, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}

	reader, err := page.PDF(&proto.PagePrintToPDF{
		PrintBackground: true,
	})
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (b *RodBrowser) GetMux() *http.ServeMux {
	return sync.OnceValue(b.NewMux)()
}

func (b RodBrowser) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/Tabs/{id}/", b.ServeHTTP)
	return mux
}

func (b *RodBrowser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.GetMux().ServeHTTP(w, r)
}

func (b *RodBrowser) ServeUi(w http.ResponseWriter, r *http.Request, upgrader *websocket.Upgrader) {
	var opts ScreencastOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		log.Println(err)
		return
	}

	page, err := b.Native().PageFromTarget(proto.TargetTargetID(r.PathValue("id")))
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	res, err := b.Native().Call(context.Background(), string(page.SessionID), "Page.enable", opts)
	if err != nil {
		log.Println(err)
		return
	}

	page.EachEvent()

	fmt.Println(string(res))

	conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
}

func (b *RodBrowser) StartScreencast(id string, opts *ScreencastOptions) (*ScreencastStream, error) {
	if opts == nil {
		opts = new(ScreencastOptions)
	}

	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return nil, err
	}

	pageWithCancel, cancel := page.WithCancel()

	frameChan := make(chan []byte)

	go func() {
		defer close(frameChan)

		pageWithCancel.EachEvent(func(e *proto.PageScreencastFrame) {
			if e == nil {
				return
			}
			frameChan <- e.Data

			if err := RodCall(page, proto.PageScreencastFrameAck{
				SessionID: e.SessionID,
			}); err != nil {

			}
		})()
	}()

	if _, err := b.Native().Call(context.Background(), string(page.SessionID), "Page.startScreencast", opts); err != nil {
		return nil, err
	}

	return NewScreencastStream(frameChan, func() error {
		cancel()

		_, err := b.Native().Call(context.Background(), string(page.SessionID), "Page.stopScreencast", opts)
		if err != nil {
			fmt.Println("failed to stop screencast", err)
		}
		return err
	}), nil
}

func GetStringsFromRodElements(elements rod.Elements, getter func(*rod.Element) (string, error)) (types.Strings, error) {
	strs, errs := stl.MapWithError(elements, true, getter)
	if err := errs.Simplify(); err != nil {
		return nil, err
	}
	return strs, nil
}

func GetTextsFromRodElements(elements rod.Elements) (types.Strings, error) {
	return GetStringsFromRodElements(elements, (*rod.Element).Text)
}

func GetHtmlsFromRodElements(elements rod.Elements) (types.Strings, error) {
	return GetStringsFromRodElements(elements, (*rod.Element).HTML)
}

type RodBrowserManager struct {
	RodBrowser
}

func RodPageToTab(page *rod.Page) *Tab {
	tab := Tab{
		Id: string(page.TargetID),
	}

	if info, err := page.Info(); err == nil {
		tab.Title = info.Title
		tab.Url = info.URL
	}

	return &tab
}

func RodPagesToTabs(pages rod.Pages) Tabs {
	return stl.Map(pages, RodPageToTab)
}

func NetworkCookieParamFromStd(cookie *http.Cookie) *proto.NetworkCookieParam {
	return &proto.NetworkCookieParam{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Domain:   cookie.Domain,
		Path:     cookie.Path,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		SameSite: NetworkCookieSameSiteFromStd(cookie.SameSite),
		Expires:  RodTimeSinceEpochFromTime(cookie.Expires),
	}
}

func NetworkCookieParamsFromStd(cookies []*http.Cookie) []*proto.NetworkCookieParam {
	return stl.Map(cookies, NetworkCookieParamFromStd)
}

func NetworkCookieToStd(cookie *proto.NetworkCookie) *http.Cookie {
	return &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Domain:   cookie.Domain,
		Path:     cookie.Path,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HTTPOnly,
		SameSite: NetworkCookieSameSiteToStd(cookie.SameSite),
		Expires:  RodTimeSinceEpochToTime(cookie.Expires),
	}
}

func NetworkCookiesToStd(cookies []*proto.NetworkCookie) []*http.Cookie {
	return stl.Map(cookies, NetworkCookieToStd)
}

func NetworkCookieSameSiteFromStd(sameSite http.SameSite) proto.NetworkCookieSameSite {
	switch sameSite {
	case http.SameSiteStrictMode:
		return proto.NetworkCookieSameSiteStrict
	case http.SameSiteLaxMode:
		return proto.NetworkCookieSameSiteLax
	case http.SameSiteNoneMode:
		return proto.NetworkCookieSameSiteNone
	}
	return proto.NetworkCookieSameSiteNone
}

func NetworkCookieSameSiteToStd(sameSite proto.NetworkCookieSameSite) http.SameSite {
	switch sameSite {
	case proto.NetworkCookieSameSiteStrict:
		return http.SameSiteStrictMode
	case proto.NetworkCookieSameSiteLax:
		return http.SameSiteLaxMode
	case proto.NetworkCookieSameSiteNone:
		return http.SameSiteNoneMode
	}
	return http.SameSiteNoneMode
}

func RodTimeSinceEpochFromTime(time time.Time) proto.TimeSinceEpoch {
	if time.IsZero() {
		return -1
	}
	return proto.TimeSinceEpoch(time.Unix())
}

func RodTimeSinceEpochToTime(timeSinceEpoch proto.TimeSinceEpoch) time.Time {
	if timeSinceEpoch < 0 {
		return time.Time{}
	}
	return timeSinceEpoch.Time()
}

func RodMouseButtonFromStd(button string) proto.InputMouseButton {
	switch button {
	case MouseButtonLeft:
		return proto.InputMouseButtonLeft
	case MouseButtonMiddle:
		return proto.InputMouseButtonMiddle
	case MouseButtonRight:
		return proto.InputMouseButtonRight
	case MouseButtonBack:
		return proto.InputMouseButtonBack
	case MouseButtonForward:
		return proto.InputMouseButtonForward
	case MouseButtonNone:
		return proto.InputMouseButtonNone
	default:
		return proto.InputMouseButtonNone
	}
}

func RodCall[
	Data interface {
		Call(c proto.Client) error
	},
](client proto.Client, data Data) (err error) {
	return data.Call(client)
}

var _ Browser = &RodBrowserManager{}
