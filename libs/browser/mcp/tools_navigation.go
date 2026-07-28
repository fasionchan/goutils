package browsermcp

import (
	"context"

	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *BrowserMcpServer) registerNavigationTools() {
	s.AddTool(mcp.NewTool(
		"browser_navigate",
		mcp.WithDescription("Navigate a tab to a URL. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL to navigate to")),
	), s.handleNavigate)

	s.AddTool(mcp.NewTool(
		"browser_navigate_back",
		mcp.WithDescription("Go back in tab history. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleNavigateBack)

	s.AddTool(mcp.NewTool(
		"browser_navigate_forward",
		mcp.WithDescription("Go forward in tab history. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleNavigateForward)

	s.AddTool(mcp.NewTool(
		"browser_reload",
		mcp.WithDescription("Reload the current page in a tab. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleReload)
}

func (s *BrowserMcpServer) handleNavigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	url, err := param.RequiredString(args, "url")
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.Navigate(tabID, url); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handleNavigateBack(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.GoBack(tabID); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handleNavigateForward(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.GoForward(tabID); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handleReload(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.Reload(tabID); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}
