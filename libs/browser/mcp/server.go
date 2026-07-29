package browsermcp

import (
	"github.com/fasionchan/goutils/libs/browser"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultMCPPath = "/mcp"
	ServerName     = "browser"
	ServerVersion  = "0.1.0"
)

// BrowserMcpServer wraps a mark3labs MCPServer bound to a Browser backend.
type BrowserMcpServer struct {
	browser browser.Browser
	path    string
	*server.MCPServer
}

type Option func(*BrowserMcpServer)

// WithPath sets the HTTP mount prefix (default /mcp).
func WithPath(path string) Option {
	return func(s *BrowserMcpServer) {
		s.path = path
	}
}

// NewBrowserMcpServer creates an MCP server and registers P0+P1 browser tools.
func NewBrowserMcpServer(b browser.Browser, opts ...Option) (*BrowserMcpServer) {
	s := &BrowserMcpServer{
		browser: b,
		path:    DefaultMCPPath,
		MCPServer: server.NewMCPServer(
			ServerName,
			ServerVersion,
			server.WithToolCapabilities(true),
		),
	}
	for _, opt := range opts {
		opt(s)
	}

	s.registerAllTools()
	return s
}

func (s *BrowserMcpServer) GetPath() string {
	if s == nil || s.path == "" {
		return DefaultMCPPath
	}
	return s.path
}

func (s *BrowserMcpServer) GetBrowser() browser.Browser {
	if s == nil {
		return nil
	}
	return s.browser
}

func (s *BrowserMcpServer) registerAllTools() {
	s.registerTabTools()
	s.registerNavigationTools()
	s.registerActionTools()
	s.registerObservationTools()
	s.registerCookieTools()
	s.registerP1Tools()
}
