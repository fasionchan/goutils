package browsermcp

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

type recordingBrowser struct {
	tabs        browser.Tabs
	lastClick   []any
	lastNavURL  string
	lastKeyEvts []*browser.KeyEvent
	getTabErr   error
	navigateErr error
}

func (b *recordingBrowser) NewTab(options *browser.NewTabOptions) (*browser.Tab, error) {
	tab := &browser.Tab{Id: "new-tab", Title: "New", Url: ""}
	if options != nil && options.Url != "" {
		tab.Url = options.Url
	}
	b.tabs = append(b.tabs, tab)
	return tab, nil
}
func (b *recordingBrowser) GetCDPAddress() (*net.TCPAddr, error) {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9222}, nil
}
func (b *recordingBrowser) GetTab(id string) (*browser.Tab, error) {
	if b.getTabErr != nil {
		return nil, b.getTabErr
	}
	for _, t := range b.tabs {
		if t != nil && t.Id == id {
			return t, nil
		}
	}
	return nil, errors.New("tab not found")
}
func (b *recordingBrowser) ListTabs() (browser.Tabs, error) { return b.tabs, nil }
func (b *recordingBrowser) SwitchToTab(id string) error     { return nil }
func (b *recordingBrowser) CloseTab(id string) error        { return nil }
func (b *recordingBrowser) Navigate(id, url string) error {
	if b.navigateErr != nil {
		return b.navigateErr
	}
	b.lastNavURL = url
	return nil
}
func (b *recordingBrowser) GoBack(id string) error    { return nil }
func (b *recordingBrowser) GoForward(id string) error { return nil }
func (b *recordingBrowser) Reload(id string) error    { return nil }
func (b *recordingBrowser) Click(id, selector, selectorType, button string, count int) error {
	b.lastClick = []any{id, selector, selectorType, button, count}
	return nil
}
func (b *recordingBrowser) Type(id, selector, selectorType, text string) error { return nil }
func (b *recordingBrowser) Hover(id, selector, selectorType string) error      { return nil }
func (b *recordingBrowser) SelectOption(id, target, targetType string, options []string, optionType string, selected bool) error {
	return nil
}
func (b *recordingBrowser) SetInputFiles(id, selector, selectorType string, files []string) error {
	return nil
}
func (b *recordingBrowser) Screenshot(id string, opts *browser.ScreenshotOptions) ([]byte, error) {
	return []byte("img"), nil
}
func (b *recordingBrowser) Snapshot(id, snapshotType string) (string, error) {
	return "role=button [ref=e1]", nil
}
func (b *recordingBrowser) GetTexts(id, target, targetType string) (types.Strings, error) {
	return types.Strings{"hello"}, nil
}
func (b *recordingBrowser) GetHtmls(id, target, targetType string) (types.Strings, error) {
	return types.Strings{"<div/>"}, nil
}
func (b *recordingBrowser) SetCookies(id string, cookies []*http.Cookie) error { return nil }
func (b *recordingBrowser) GetCookies(id string) ([]*http.Cookie, error) {
	return []*http.Cookie{{Name: "a", Value: "b"}}, nil
}
func (b *recordingBrowser) PrintToPdf(id string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("%PDF")), nil
}
func (b *recordingBrowser) StartScreencast(id string, opts *browser.ScreencastOptions) (*browser.ScreencastStream, error) {
	return nil, nil
}
func (b *recordingBrowser) Close() error { return nil }
func (b *recordingBrowser) DispatchMouseEvent(id string, event *browser.MouseEvent) error {
	return nil
}
func (b *recordingBrowser) DispatchKeyEvent(id string, event *browser.KeyEvent) error {
	copied := *event
	b.lastKeyEvts = append(b.lastKeyEvts, &copied)
	return nil
}
func (b *recordingBrowser) GetScreencastSessionMeta(id string, opts *browser.ScreencastOptions) (*browser.ScreencastSessionMeta, error) {
	return &browser.ScreencastSessionMeta{}, nil
}

func callTool(t *testing.T, s *BrowserMcpServer, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	tools := s.ListTools()
	tool, ok := tools[name]
	require.True(t, ok, "tool %s missing", name)
	require.NotNil(t, tool.Handler)
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	result, err := tool.Handler(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)
	return result
}

func TestNewBrowserMcpServerNilBrowser(t *testing.T) {
	s := NewBrowserMcpServer(nil)
	require.NotNil(t, s)
}

func TestMissingTabID(t *testing.T) {
	s := NewBrowserMcpServer(&recordingBrowser{})
	require.NotNil(t, s)

	result := callTool(t, s, "browser_navigate", map[string]any{"url": "https://example.com"})
	require.True(t, result.IsError)
	require.Contains(t, result.Content[0].(mcp.TextContent).Text, "tabId is required")
}

func TestNavigateAndClickMapping(t *testing.T) {
	b := &recordingBrowser{tabs: browser.Tabs{{Id: "t1", Title: "x", Url: "about:blank"}}}
	s := NewBrowserMcpServer(b)
	require.NotNil(t, s)

	result := callTool(t, s, "browser_navigate", map[string]any{
		"tabId": "t1",
		"url":   "https://example.com",
	})
	require.False(t, result.IsError)
	require.Equal(t, "https://example.com", b.lastNavURL)

	result = callTool(t, s, "browser_click", map[string]any{
		"tabId":       "t1",
		"target":      "e1",
		"targetType":  "ref",
		"doubleClick": true,
	})
	require.False(t, result.IsError)
	require.Equal(t, []any{"t1", "e1", "ref", "left", 2}, b.lastClick)
}

func TestUnknownTabError(t *testing.T) {
	b := &recordingBrowser{getTabErr: errors.New("tab not found")}
	s := NewBrowserMcpServer(b)
	require.NotNil(t, s)

	result := callTool(t, s, "browser_get_tab", map[string]any{"tabId": "missing"})
	require.True(t, result.IsError)
	require.Contains(t, result.Content[0].(mcp.TextContent).Text, "tab not found")
}

func TestToolListAudit(t *testing.T) {
	s := NewBrowserMcpServer(&recordingBrowser{})
	require.NotNil(t, s)

	required := []string{
		"browser_list_tabs", "browser_new_tab", "browser_close_tab", "browser_get_tab",
		"browser_select_tab",
		"browser_navigate", "browser_navigate_back", "browser_navigate_forward", "browser_reload",
		"browser_click", "browser_type", "browser_hover", "browser_select_option",
		"browser_snapshot", "browser_take_screenshot", "browser_get_texts", "browser_get_htmls",
		"browser_cookie_list", "browser_cookie_set",
		"browser_file_upload", "browser_press_key", "browser_pdf_save",
	}
	tools := s.ListTools()
	for _, name := range required {
		_, ok := tools[name]
		require.True(t, ok, "missing tool %s", name)
	}

	forbidden := []string{
		"browser_wait_for", "browser_evaluate", "browser_handle_dialog",
		"browser_network_requests", "browser_console_messages",
		"browser_remote", "browser_screencast",
	}
	for _, name := range forbidden {
		_, ok := tools[name]
		require.False(t, ok, "unexpected tool %s", name)
	}
}

func TestPressKeyDispatchesDownUp(t *testing.T) {
	b := &recordingBrowser{}
	s := NewBrowserMcpServer(b)
	require.NotNil(t, s)

	result := callTool(t, s, "browser_press_key", map[string]any{
		"tabId": "t1",
		"key":   "Enter",
	})
	require.False(t, result.IsError)
	require.Len(t, b.lastKeyEvts, 2)
	require.Equal(t, browser.KeyEventTypeDown, b.lastKeyEvts[0].Type)
	require.Equal(t, browser.KeyEventTypeUp, b.lastKeyEvts[1].Type)
}

func TestHTTPMountPaths(t *testing.T) {
	s := NewBrowserMcpServer(&recordingBrowser{}, WithPath("/mcp"))
	require.NotNil(t, s)
	h := s.MountOnto(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	// non-mcp path falls through
	req, _ := http.NewRequest(http.MethodGet, "/Tabs", nil)
	rr := newRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusTeapot, rr.code)
}

type recorder struct {
	code int
	hdr  http.Header
	body strings.Builder
}

func newRecorder() *recorder {
	return &recorder{code: 200, hdr: make(http.Header)}
}

func (r *recorder) Header() http.Header         { return r.hdr }
func (r *recorder) Write(b []byte) (int, error) { return r.body.Write(b) }
func (r *recorder) WriteHeader(statusCode int)  { r.code = statusCode }
