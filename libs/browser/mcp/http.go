package browsermcp

import (
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

const (
	streamableSuffix = "/streamable"
	sseSuffix        = "/sse"
	sseMessageSuffix = "/sse/message"
)

// NewHTTPHandler returns an http.Handler serving Streamable HTTP and SSE under the MCP path prefix.
func (s *BrowserMcpServer) NewHTTPHandler() http.Handler {
	prefix := ensurePathPrefix(s.GetPath())
	mux := http.NewServeMux()

	streamable := server.NewStreamableHTTPServer(s.MCPServer)
	mux.Handle(prefix, streamable)
	mux.Handle(prefix+streamableSuffix, streamable)

	sse := server.NewSSEServer(
		s.MCPServer,
		server.WithSSEEndpoint(prefix+sseSuffix),
		server.WithMessageEndpoint(prefix+sseMessageSuffix),
	)
	mux.Handle(prefix+sseSuffix, sse)
	mux.Handle(prefix+sseMessageSuffix, sse)

	return mux
}

// MountOnto combines an existing API handler with the MCP HTTP endpoints.
// Requests under the MCP path go to MCP; everything else goes to apiHandler.
func (s *BrowserMcpServer) MountOnto(apiHandler http.Handler) http.Handler {
	if apiHandler == nil {
		apiHandler = http.NotFoundHandler()
	}
	mcpHandler := s.NewHTTPHandler()
	prefix := ensurePathPrefix(s.GetPath())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			mcpHandler.ServeHTTP(w, r)
			return
		}
		apiHandler.ServeHTTP(w, r)
	})
}
