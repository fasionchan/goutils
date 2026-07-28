package browsermcp

import (
	"context"
	"net/http"

	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *BrowserMcpServer) registerCookieTools() {
	s.AddTool(mcp.NewTool(
		"browser_cookie_list",
		mcp.WithDescription("List cookies for a tab. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleCookieList)

	s.AddTool(mcp.NewTool(
		"browser_cookie_set",
		mcp.WithDescription("Set a cookie on a tab (flat fields). Requires tabId, name, value."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Cookie name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Cookie value")),
		mcp.WithString("domain", mcp.Description("Cookie domain")),
		mcp.WithString("path", mcp.Description("Cookie path")),
		mcp.WithBoolean("secure", mcp.Description("Secure flag"), mcp.DefaultBool(false)),
		mcp.WithBoolean("httpOnly", mcp.Description("HttpOnly flag"), mcp.DefaultBool(false)),
	), s.handleCookieSet)
}

func (s *BrowserMcpServer) handleCookieList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	cookies, err := s.browser.GetCookies(tabID)
	if err != nil {
		return toolError(err)
	}
	return toolJSON(cookies)
}

func (s *BrowserMcpServer) handleCookieSet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	cookie, err := cookieFromArgs(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.SetCookies(tabID, []*http.Cookie{cookie}); err != nil {
		return toolError(err)
	}
	return toolJSON(cookie)
}
