package browsermcp

import (
	"context"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *BrowserMcpServer) registerTabTools() {
	s.AddTool(mcp.NewTool(
		"browser_list_tabs",
		mcp.WithDescription("List all open browser tabs (id, title, url)."),
	), s.handleListTabs)

	s.AddTool(mcp.NewTool(
		"browser_new_tab",
		mcp.WithDescription("Create a new browser tab. Optionally navigate to url."),
		mcp.WithString("url", mcp.Description("Optional URL to open in the new tab")),
		mcp.WithNumber("width", mcp.Description("Optional viewport width")),
		mcp.WithNumber("height", mcp.Description("Optional viewport height")),
	), s.handleNewTab)

	s.AddTool(mcp.NewTool(
		"browser_close_tab",
		mcp.WithDescription("Close a browser tab by tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleCloseTab)

	s.AddTool(mcp.NewTool(
		"browser_get_tab",
		mcp.WithDescription("Get tab metadata (id, title, url) by tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleGetTab)
}

func (s *BrowserMcpServer) handleListTabs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tabs, err := s.browser.ListTabs()
	if err != nil {
		return toolError(err)
	}
	return toolJSON(tabs)
}

func (s *BrowserMcpServer) handleNewTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	opts := browser.NewNewTabOptions()
	url, err := param.OptionalString(args, "url")
	if err != nil {
		return toolError(err)
	}
	if url != "" {
		opts.WithUrl(url)
	}
	if param.Has(args, "width") {
		w, err := param.OptionalInt(args, "width", 0)
		if err != nil {
			return toolError(err)
		}
		opts.WithWidth(w)
	}
	if param.Has(args, "height") {
		h, err := param.OptionalInt(args, "height", 0)
		if err != nil {
			return toolError(err)
		}
		opts.WithHeight(h)
	}
	tab, err := s.browser.NewTab(opts)
	if err != nil {
		return toolError(err)
	}
	return toolJSON(tab)
}

func (s *BrowserMcpServer) handleCloseTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.CloseTab(tabID); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handleGetTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	tab, err := s.browser.GetTab(tabID)
	if err != nil {
		return toolError(err)
	}
	return toolJSON(tab)
}
