package browser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

const (
	SnapshotTypeA11y = "a11y"
	SnapshotTypeDom  = "dom"

	LocatorTypeCssSelector = "css-selector"
	LocatorTypeXPath       = "xpath"
	LocatorTypeRef         = "ref"

	OptionLocatorTypeText        = "text"
	OptionLocatorTypeCssSelector = "css-selector"
	OptionLocatorTypeRegex       = "regex"

	MouseButtonLeft    = "left"
	MouseButtonMiddle  = "middle"
	MouseButtonRight   = "right"
	MouseButtonBack    = "back"
	MouseButtonForward = "forward"
	MouseButtonNone    = "none"

	MouseEventTypeMove  = "move"
	MouseEventTypeDown  = "down"
	MouseEventTypeUp    = "up"
	MouseEventTypeWheel = "wheel"

	KeyEventTypeDown = "down"
	KeyEventTypeUp   = "up"

	InputModifierAlt   = "alt"
	InputModifierCtrl  = "ctrl"
	InputModifierMeta  = "meta"
	InputModifierShift = "shift"
)

type Browser interface {
	GetCDPAddress() (*net.TCPAddr, error)

	NewTab(options *NewTabOptions) (*Tab, error)
	GetTab(id string) (*Tab, error)
	ListTabs() (Tabs, error)
	SwitchToTab(id string) error
	CloseTab(id string) error

	Navigate(id, url string) error
	GoBack(id string) error
	GoForward(id string) error
	Reload(id string) error

	Click(id, selector, selectorType, button string, count int) error
	Type(id, selector, selectorType, text string) error
	Hover(id, selector, selectorType string) error
	SelectOption(id, target, targetType string, options []string, optionType string, selected bool) error
	SetInputFiles(id, selector, selectorType string, files []string) error

	DispatchMouseEvent(id string, event *MouseEvent) error
	DispatchKeyEvent(id string, event *KeyEvent) error

	Screenshot(id string, opts *ScreenshotOptions) ([]byte, error)
	Snapshot(id, snapshotType string) (string, error)
	GetTexts(id, target, targetType string) (types.Strings, error)
	GetHtmls(id, target, targetType string) (types.Strings, error)

	SetCookies(id string, cookies []*http.Cookie) error
	GetCookies(id string) ([]*http.Cookie, error)

	PrintToPdf(id string) (io.ReadCloser, error)

	StartScreencast(id string, opts *ScreencastOptions) (*ScreencastStream, error)
	GetScreencastSessionMeta(id string, opts *ScreencastOptions) (*ScreencastSessionMeta, error)

	Close() error
}

// MouseEvent is a coordinate-based mouse input in page CSS pixels.
type MouseEvent struct {
	Type       string   `json:"type"`
	X          float64  `json:"x"`
	Y          float64  `json:"y"`
	Button     string   `json:"button,omitempty"`
	ClickCount int      `json:"click_count,omitempty"`
	DeltaX     float64  `json:"delta_x,omitempty"`
	DeltaY     float64  `json:"delta_y,omitempty"`
	Modifiers  []string `json:"modifiers,omitempty"`
}

// KeyEvent is a keyboard input event.
type KeyEvent struct {
	Type       string   `json:"type"`
	Key        string   `json:"key,omitempty"`
	Code       string   `json:"code,omitempty"`
	Text       string   `json:"text,omitempty"`
	Modifiers  []string `json:"modifiers,omitempty"`
	AutoRepeat bool     `json:"auto_repeat,omitempty"`
}

// ScreencastSessionMeta describes viewport/frame geometry for remote control clients.
type ScreencastSessionMeta struct {
	Format            string  `json:"format,omitempty"`
	ViewportWidth     int     `json:"viewport_width"`
	ViewportHeight    int     `json:"viewport_height"`
	FrameWidth        int     `json:"frame_width"`
	FrameHeight       int     `json:"frame_height"`
	DeviceScaleFactor float64 `json:"device_scale_factor,omitempty"`
}

// EncodeInputModifiers converts modifier names to CDP bit flags (Alt=1, Ctrl=2, Meta=4, Shift=8).
func EncodeInputModifiers(modifiers []string) int {
	var bits int
	for _, modifier := range modifiers {
		switch strings.ToLower(strings.TrimSpace(modifier)) {
		case InputModifierAlt:
			bits |= 1
		case InputModifierCtrl, "control":
			bits |= 2
		case InputModifierMeta, "command", "cmd":
			bits |= 4
		case InputModifierShift:
			bits |= 8
		}
	}
	return bits
}

func EstimateScreencastFrameSize(viewportWidth, viewportHeight int, opts *ScreencastOptions) (int, int) {
	width, height := viewportWidth, viewportHeight
	if width <= 0 || height <= 0 {
		return width, height
	}
	if opts == nil {
		return width, height
	}

	scale := 1.0
	if opts.MaxWidth != nil && *opts.MaxWidth > 0 {
		if s := float64(*opts.MaxWidth) / float64(width); s < scale {
			scale = s
		}
	}
	if opts.MaxHeight != nil && *opts.MaxHeight > 0 {
		if s := float64(*opts.MaxHeight) / float64(height); s < scale {
			scale = s
		}
	}
	if scale >= 1 {
		return width, height
	}
	return int(float64(width) * scale), int(float64(height) * scale)
}

func (opts *ScreencastOptions) GetFormat() string {
	if opts == nil || opts.Format == nil || *opts.Format == "" {
		return "jpeg"
	}
	return *opts.Format
}

type NewTabOptions struct {
	Url    string
	Width  *int
	Height *int
}

func NewNewTabOptions(options ...NewTabOption) *NewTabOptions {
	return stl.NewOptions(options...).Apply(new(NewTabOptions))
}

func (opts *NewTabOptions) WithUrl(url string) *NewTabOptions {
	opts.Url = url
	return opts
}

func (opts *NewTabOptions) WithWidth(width int) *NewTabOptions {
	opts.Width = &width
	return opts
}

func (opts *NewTabOptions) WithHeight(height int) *NewTabOptions {
	opts.Height = &height
	return opts
}

type NewTabOption = stl.Option[*NewTabOptions]

func NewTabWithUrl(url string) NewTabOption {
	return func(opts *NewTabOptions) {
		opts.Url = url
	}
}

func NewTabWithWidth(width int) NewTabOption {
	return func(opts *NewTabOptions) {
		opts.Width = &width
	}
}

func NewTabWithHeight(height int) NewTabOption {
	return func(opts *NewTabOptions) {
		opts.Height = &height
	}
}

type TabPtr = *Tab

type Tab struct {
	Id    string
	Title string
	Url   string
}

func (tab *Tab) GetId() string {
	return tab.Id
}

func (tab *Tab) GetTitle() string {
	return tab.Title
}

func (tab *Tab) GetUrl() string {
	return tab.Url
}

type Tabs []*Tab

func (tabs Tabs) Ids() types.Strings {
	return stl.Map(tabs, TabPtr.GetId)
}

type ViewRegion struct {
}

type Viewport struct {
	// X offset in device independent pixels (dip).
	X float64 `json:"x" query:"x"`

	// Y offset in device independent pixels (dip).
	Y float64 `json:"y" query:"y"`

	// Width Rectangle width in device independent pixels (dip).
	Width float64 `json:"width" query:"width"`

	// Height Rectangle height in device independent pixels (dip).
	Height float64 `json:"height" query:"height"`
}

type ScreenshotOptions struct {
	Format  *string    `json:"format,omitempty"`
	Quality *int       `json:"quality,omitempty"`
	Clip    *Viewport  `json:"clip,omitempty"`
	Target  *TypedExpr `json:"target,omitempty"`

	// Scale Page scale factor.
	Scale *float64 `json:"scale,omitempty"`
}

type ScreenOptionsQuery struct {
	Format  *string    `query:"format" required:"true" default:"jpeg" enum:"jpeg,png"`
	Quality *int       `query:"quality" required:"true" default:"60" min:"0" max:"100"`
	Clip    *Viewport  `query:"clip"`
	Target  *TypedExpr `query:"target"`

	// Scale Page scale factor.
	Scale *float64 `query:"scale,omitempty"`
}

func (query *ScreenOptionsQuery) ToScreenshotOptions() *ScreenshotOptions {
	return &ScreenshotOptions{
		Format:  query.Format,
		Quality: query.Quality,
		Clip:    query.Clip,
		Target:  query.Target,
		Scale:   query.Scale,
	}
}

type TypedExpr struct {
	Expr string `json:"expr" query:"expr"`
	Type string `json:"type" query:"type"`
}

func NewScreenshotOptions(options ...ScreenshotOption) *ScreenshotOptions {
	return stl.NewOptions(options...).Apply(new(ScreenshotOptions))
}

func NewScreenshotOptionsFromUrlValues(query url.Values) (*ScreenshotOptions, error) {
	var opts stl.Options[*ScreenshotOptions]

	if format := query.Get("format"); format != "" {
		opts = append(opts, ScreenshotWithFormat(format))
	}

	opts, err := opts.ParseAndAppendIntOptions(stl.IntOptionMetas[*ScreenshotOptions]{
		{Opt: ScreenshotWithQuality, Value: query.Get("quality")},
	})
	if err != nil {
		return nil, err
	}

	return NewScreenshotOptions(opts...), nil

}

func (opts *ScreenshotOptions) GetFormat() string {
	if opts == nil {
		return "png"
	}

	if format := opts.Format; format != nil {
		return *format
	}

	return "png"
}

func (opts *ScreenshotOptions) MimeType() string {
	return fmt.Sprintf("image/%s", opts.GetFormat())
}

type ScreenshotOption = stl.Option[*ScreenshotOptions]

func ScreenshotWithFormat(format string) ScreenshotOption {
	return func(opts *ScreenshotOptions) {
		opts.Format = &format
	}
}

func ScreenshotWithQuality(quality int) ScreenshotOption {
	return func(opts *ScreenshotOptions) {
		opts.Quality = &quality
	}
}

type ScreencastOptions struct {
	Format        *string `json:"format,omitempty"`
	Quality       *int    `json:"quality,omitempty"`
	MaxWidth      *int    `json:"max_width,omitempty"`
	MaxHeight     *int    `json:"max_height,omitempty"`
	EventNthFrame *int    `json:"event_nth_frame,omitempty"`
}

func NewScreencastOptions(options ...ScreencastOption) *ScreencastOptions {
	return stl.NewOptions(options...).Apply(new(ScreencastOptions))
}

func NewScreencastOptionsFromUrlValues(query url.Values) (*ScreencastOptions, error) {
	var opts stl.Options[*ScreencastOptions]

	if format := query.Get("format"); format != "" {
		opts = append(opts, ScreencastWithFormat(format))
	}

	opts, err := opts.ParseAndAppendIntOptions(stl.IntOptionMetas[*ScreencastOptions]{
		{Opt: ScreencastWithQuality, Value: query.Get("quality")},
		{Opt: ScreencastWithMaxWidth, Value: query.Get("max_width")},
		{Opt: ScreencastWithMaxHeight, Value: query.Get("max_height")},
		{Opt: ScreencastWithEventNthFrame, Value: query.Get("event_nth_frame")},
	})
	if err != nil {
		return nil, err
	}

	return NewScreencastOptions(opts...), nil
}

type ScreencastOption = stl.Option[*ScreencastOptions]

func ScreencastWithFormat(format string) ScreencastOption {
	return func(opts *ScreencastOptions) {
		opts.Format = &format
	}
}

func ScreencastWithQuality(quality int) ScreencastOption {
	return func(opts *ScreencastOptions) {
		opts.Quality = &quality
	}
}

func ScreencastWithMaxWidth(maxWidth int) ScreencastOption {
	return func(opts *ScreencastOptions) {
		opts.MaxWidth = &maxWidth
	}
}

func ScreencastWithMaxHeight(maxHeight int) ScreencastOption {
	return func(opts *ScreencastOptions) {
		opts.MaxHeight = &maxHeight
	}
}

func ScreencastWithEventNthFrame(eventNthFrame int) ScreencastOption {
	return func(opts *ScreencastOptions) {
		opts.EventNthFrame = &eventNthFrame
	}
}

type BytesChan chan []byte

type CloseFunc func() error

func (fn CloseFunc) Close() error {
	return fn()
}

type ScreencastStream struct {
	BytesChan
	CloseFunc
}

func NewScreencastStream(frameChan BytesChan, closeFunc CloseFunc) *ScreencastStream {
	return &ScreencastStream{
		BytesChan: frameChan,
		CloseFunc: closeFunc,
	}
}

type BrowserLaunchOptions struct {
	Addr  *net.TCPAddr
	Flags url.Values
}

func NewBrowserLaunchOptionsFromEnv(getenv func(string) string) *BrowserLaunchOptions {
	addrStr := getenv("BROWSER_ADDR")
	if addrStr == "" {
		addrStr = "127.0.0.1:9222"
	}

	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil
	}

	flags := url.Values{}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if !strings.HasPrefix(name, "BROWSER_") {
			continue
		}

		name = strings.TrimPrefix(name, "BROWSER_")
		name = strings.ToLower(name)
		name = strings.ReplaceAll(name, "_", "-")

		value := strings.TrimSpace(parts[1])
		flags.Add(name, value)
	}

	return &BrowserLaunchOptions{
		Addr:  addr,
		Flags: flags,
	}
}

// deep copy?
func (opts *BrowserLaunchOptions) Dup() *BrowserLaunchOptions {
	return stl.Dup(opts)
}

func (opts *BrowserLaunchOptions) WithAddr(addr *net.TCPAddr) *BrowserLaunchOptions {
	if opts == nil {
		return nil
	}

	opts.Addr = addr
	return opts
}

func (opts *BrowserLaunchOptions) WithFlag(name, value string) *BrowserLaunchOptions {
	if opts == nil {
		return nil
	}

	if flags := opts.Flags; flags != nil {
		flags.Add(name, value)
	}

	return opts
}

type BrowserLauncher interface {
	Launch(ctx context.Context, opts *BrowserLaunchOptions) (Browser, error)
}

type BrowserLaunchFunc func(ctx context.Context, opts *BrowserLaunchOptions) (Browser, error)

func (fn BrowserLaunchFunc) Launch(ctx context.Context, opts *BrowserLaunchOptions) (Browser, error) {
	if fn == nil {
		return nil, errors.New("browser launch function is nil")
	}
	return fn(ctx, opts)
}
