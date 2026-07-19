package browser

import (
	"io"
	"net/http"
	"net/url"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

const (
	SnapshotTypeA11y = "a11y"
	SnapshotTypeDom  = "dom"

	SelectorTypeCss   = "css"
	SelectorTypeXPath = "xpath"
	SelectorTypeRef   = "ref"

	MouseButtonLeft    = "left"
	MouseButtonMiddle  = "middle"
	MouseButtonRight   = "right"
	MouseButtonBack    = "back"
	MouseButtonForward = "forward"
	MouseButtonNone    = "none"
)

type Browser interface {
	NewTab(options *NewTabOptions) (*Tab, error)
	ListTabs() (Tabs, error)
	SwitchToTab(id string) error
	CloseTab(id string) error

	Navigate(id, url string) error
	GoBack(id string) error
	GoForward(id string) error
	Reload(id string) error

	Click(id, selector, selectorType, button string, count int) error
	Type(id, selector, selectorType, text string) error
	SetInputFiles(id, selector, selectorType string, files []string) error

	Screenshot(id string) ([]byte, error)
	// Snapshot(id, snapshotType string) (string, error)
	GetTexts(id, selector, selectorType string) (types.Strings, error)
	GetHtmls(id, selector, selectorType string) (types.Strings, error)

	SetCookies(id string, cookies []*http.Cookie) error
	GetCookies(id string) ([]*http.Cookie, error)

	PrintToPdf(id string) (io.ReadCloser, error)

	StartScreencast(id string, opts *ScreencastOptions) (*ScreencastStream, error)
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

	opts, err := opts.ParseAndAppendIntOptions(stl.KeyValuePairPtrs[func(int) ScreencastOption, string]{
		{Key: ScreencastWithQuality, Value: query.Get("quality")},
		{Key: ScreencastWithMaxWidth, Value: query.Get("max_width")},
		{Key: ScreencastWithMaxHeight, Value: query.Get("max_height")},
		{Key: ScreencastWithEventNthFrame, Value: query.Get("event_nth_frame")},
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
