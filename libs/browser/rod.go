package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
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

func ConnectRodBrowserForManager() (*RodBrowserManager, error) {
	browser, err := ConnectRodBrowser()
	if err != nil {
		return nil, err
	}
	return NewRodBrowserManagerFromNative(browser.Native()), nil
}

func LaunchRodBrowser(ctx context.Context, opts *BrowserLaunchOptions) (*RodBrowser, error) {
	launcher := launcher.New()

	launcher.Set("headless", "new")
	// launcher.Set("disable-gpu", "true")

	if addr := opts.Addr; addr != nil {
		launcher.Set("remote-debugging-address", addr.IP.String())
		launcher.Set("remote-debugging-port", strconv.Itoa(addr.Port))
	}

	for key, values := range opts.Flags {
		launcher.Set(flags.Flag(key), values...)
	}

	url, err := launcher.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(url)
	if err := browser.Connect(); err != nil {
		return nil, err
	}

	return (*RodBrowser)(browser), nil
}

func LaunchRodBrowserForManager(ctx context.Context, opts *BrowserLaunchOptions) (*RodBrowserManager, error) {
	browser, err := LaunchRodBrowser(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewRodBrowserManagerFromNative(browser.Native()), nil
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

	if target := opts.Target; target != nil && opts.Clip == nil {
		elements, err := b.getElementsFromPage(page, target.Expr, target.Type)
		if err != nil {
			return nil, err
		}

		if len(elements) == 0 {
			return nil, fmt.Errorf("no elements found for selector: %s", target.Expr)
		}

		results, err := RodElements(elements).Shape()
		if err != nil {
			return nil, err
		}

		rect := results.Box().Box()
		if rect == nil {
			return nil, fmt.Errorf("no rect found for selector: %s", target.Expr)
		}

		opts.Clip = &Viewport{
			X:      rect.X,
			Y:      rect.Y,
			Width:  rect.Width,
			Height: rect.Height,
		}
	}

	clip := PageCaptureScreenshotClipFromOptions(opts.Clip)
	if scale := opts.Scale; scale != nil {
		clip.Scale = *scale
	}

	return page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormat(opts.GetFormat()),
		Quality: opts.Quality,
		Clip:    clip,
	})
}

func PageCaptureScreenshotFromOptions(opts *ScreenshotOptions) *proto.PageCaptureScreenshot {
	clip := PageCaptureScreenshotClipFromOptions(opts.Clip)
	if scale := opts.Scale; scale != nil {
		clip.Scale = *scale
	}

	return &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormat(opts.GetFormat()),
		Quality: opts.Quality,
		Clip:    clip,
	}
}

func PageCaptureScreenshotClipFromOptions(clip *Viewport) *proto.PageViewport {
	if clip == nil {
		return nil
	}

	return &proto.PageViewport{
		X:      clip.X,
		Y:      clip.Y,
		Width:  clip.Width,
		Height: clip.Height,
		Scale:  1,
	}
}

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

	return b.getElementsFromPage(page, selector, selectorType)
}

func (b *RodBrowser) getElementsFromPage(page *rod.Page, selector, selectorType string) (rod.Elements, error) {
	var fn func(*rod.Page, string) (rod.Elements, error)

	switch selectorType {
	case SelectorTypeCss:
		fn = (*rod.Page).Elements
	case SelectorTypeXPath:
		fn = (*rod.Page).ElementsX
	case SelectorTypeRef:
		fn = b.getElementsFromPageByRef
	default:
		return nil, fmt.Errorf("invalid selector type: %s", selectorType)
	}

	return fn(page, selector)
}

func (b *RodBrowser) getElementsFromPageByRef(page *rod.Page, ref string) (rod.Elements, error) {
	nodeId, err := strconv.ParseInt(ref, 10, 64)
	if err != nil {
		return nil, err
	}

	request := proto.DOMDescribeNode{
		BackendNodeID: proto.DOMBackendNodeID(nodeId),
	}

	result, err := request.Call(page)
	if err != nil {
		return nil, err
	}

	element, err := page.ElementFromNode(result.Node)
	if err != nil {
		return nil, err
	}

	return rod.Elements{element}, nil
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
	*RodBrowser
	refs *stl.SyncMap[string, IdMappingByRef]
}

func NewRodBrowserManager(b *RodBrowser) *RodBrowserManager {
	return &RodBrowserManager{
		RodBrowser: b,
		refs: stl.NewSyncMap[string, IdMappingByRef](),
	}
}

func NewRodBrowserManagerFromNative(b *rod.Browser) *RodBrowserManager {
	return NewRodBrowserManager((*RodBrowser)(b))
}

func (b *RodBrowserManager) Snapshot(id, snapshotType string) (string, error) {
	page, err := b.Native().PageFromTarget(proto.TargetTargetID(id))
	if err != nil {
		return "", err
	}

	switch snapshotType {
	case "", SnapshotTypeA11y:
		var request proto.AccessibilityGetFullAXTree
		result, err := request.Call(page)
		if err != nil {
			return "", err
		}
		return AccessibilityAxNodes(result.Nodes).String(), nil
	default:
		return "", fmt.Errorf("unsupported snapshot type: %s", snapshotType)
	}
}

type IdMappingByRef = stl.Mapping[string, string]

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

type RodElements rod.Elements

func (els RodElements) Native() rod.Elements {
	return rod.Elements(els)
}

func (els RodElements) Describe(depth int, pierce bool) (RodDomNodes, error) {
	nodes, errs := stl.MapWithError(els, true, func(el *rod.Element) (*proto.DOMNode, error) {
		return el.Describe(depth, pierce)
	})
	if err := errs.Simplify(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (els RodElements) Shape() (RodContentQuadsResults, error) {
	results, errs := stl.MapWithError(els, true, (*rod.Element).Shape)
	if err := errs.Simplify(); err != nil {
		return nil, err
	}
	return results, nil
}

type RodDomNodes []*proto.DOMNode

type RodContentQuadsResults []*proto.DOMGetContentQuadsResult

func (results RodContentQuadsResults) Box() DomRects {
	return stl.Map(results, (*proto.DOMGetContentQuadsResult).Box)
}

type DomRects []*proto.DOMRect

func (rects DomRects) Empty() bool {
	return len(rects) == 0
}

func (rects DomRects) PurgeNil() DomRects {
	return stl.PurgeZero(rects)
}

func (rects DomRects) Box() *proto.DOMRect {
	rects = rects.PurgeNil()

	if rects.Empty() {
		return nil
	}

	first := rects[0]

	x1 := first.X
	y1 := first.Y
	x2 := first.X + first.Width
	y2 := first.Y + first.Height

	for _, rect := range rects[1:] {
		x1 = min(x1, rect.X)
		y1 = min(y1, rect.Y)
		x2 = max(x2, rect.X+rect.Width)
		y2 = max(y2, rect.Y+rect.Height)
	}

	return &proto.DOMRect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}
}

var _ Browser = &RodBrowserManager{}
