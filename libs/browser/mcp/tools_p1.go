package browsermcp

import (
	"context"
	"io"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *BrowserMcpServer) registerP1Tools() {
	s.AddTool(mcp.NewTool(
		"browser_select_tab",
		mcp.WithDescription("Activate/switch to a tab by tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handleSelectTab)

	uploadOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Set input files on a file input element. Requires tabId, target, targetType, paths."),
	}, locatorToolOptions(
		mcp.WithArray("paths", mcp.Required(), mcp.Description("Absolute file paths to upload"), mcp.WithStringItems()),
	)...)
	s.AddTool(mcp.NewTool("browser_file_upload", uploadOpts...), s.handleFileUpload)

	s.AddTool(mcp.NewTool(
		"browser_press_key",
		mcp.WithDescription("Press a key (keydown+keyup) on a tab. Requires tabId and key."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Key name or character, e.g. Enter or a")),
		mcp.WithArray("modifiers", mcp.Description("Modifier keys: alt, ctrl, meta, shift"), mcp.WithStringItems()),
	), s.handlePressKey)

	s.AddTool(mcp.NewTool(
		"browser_pdf_save",
		mcp.WithDescription("Print the tab page to PDF. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
	), s.handlePDFSave)
}

func (s *BrowserMcpServer) handleSelectTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.SwitchToTab(tabID); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handleFileUpload(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	paths, err := param.RequiredStringSlice(args, "paths")
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.SetInputFiles(loc.TabID, loc.Target, loc.TargetType, paths); err != nil {
		return toolError(err)
	}
	return toolOK(loc.TabID)
}

func (s *BrowserMcpServer) handlePressKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	key, err := param.RequiredString(args, "key")
	if err != nil {
		return toolError(err)
	}
	modifiers, err := param.OptionalStringSlice(args, "modifiers")
	if err != nil {
		return toolError(err)
	}
	down := &browser.KeyEvent{Type: browser.KeyEventTypeDown, Key: key, Modifiers: modifiers}
	up := &browser.KeyEvent{Type: browser.KeyEventTypeUp, Key: key, Modifiers: modifiers}
	if err := s.browser.DispatchKeyEvent(tabID, down); err != nil {
		return toolError(err)
	}
	if err := s.browser.DispatchKeyEvent(tabID, up); err != nil {
		return toolError(err)
	}
	return toolOK(tabID)
}

func (s *BrowserMcpServer) handlePDFSave(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	reader, err := s.browser.PrintToPdf(tabID)
	if err != nil {
		return toolError(err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return toolError(err)
	}
	return toolPDF(data)
}
