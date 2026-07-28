package browsermcp

import (
	"context"
	"fmt"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

const snapshotRefHint = "\n\n---\nTip: use refs from this snapshot with targetType=ref for click/type/hover/select_option."

func (s *BrowserMcpServer) registerObservationTools() {
	s.AddTool(mcp.NewTool(
		"browser_snapshot",
		mcp.WithDescription("Capture page snapshot (default a11y). Prefer this over screenshot for actions; use returned refs with targetType=ref. Requires tabId."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("type", mcp.Description("Snapshot type"), mcp.Enum(snapshotTypeAllowed...), mcp.DefaultString(browser.SnapshotTypeA11y)),
	), s.handleSnapshot)

	s.AddTool(mcp.NewTool(
		"browser_take_screenshot",
		mcp.WithDescription("Take a screenshot of the tab. Prefer browser_snapshot for actions. Requires tabId. Default format png."),
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("type", mcp.Description("Image format"), mcp.Enum("png", "jpeg"), mcp.DefaultString("png")),
		mcp.WithNumber("quality", mcp.Description("JPEG quality 0-100"), mcp.DefaultNumber(80)),
	), s.handleTakeScreenshot)

	getTextsOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Get text contents of matching elements. Requires tabId, target, targetType."),
	}, locatorToolOptions()...)
	s.AddTool(mcp.NewTool("browser_get_texts", getTextsOpts...), s.handleGetTexts)

	getHtmlsOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Get HTML of matching elements. Requires tabId, target, targetType."),
	}, locatorToolOptions()...)
	s.AddTool(mcp.NewTool("browser_get_htmls", getHtmlsOpts...), s.handleGetHtmls)
}

func (s *BrowserMcpServer) handleSnapshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	snapType, err := param.OptionalStringEnum(args, "type", snapshotTypeAllowed, browser.SnapshotTypeA11y)
	if err != nil {
		return toolError(err)
	}
	snapshot, err := s.browser.Snapshot(tabID, snapType)
	if err != nil {
		return toolError(err)
	}
	return toolText(snapshot + snapshotRefHint)
}

func (s *BrowserMcpServer) handleTakeScreenshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	tabID, err := parseTabID(args)
	if err != nil {
		return toolError(err)
	}
	format, err := param.OptionalStringEnum(args, "type", []string{"png", "jpeg"}, "png")
	if err != nil {
		return toolError(err)
	}
	quality, err := param.OptionalInt(args, "quality", 80)
	if err != nil {
		return toolError(err)
	}
	opts := browser.NewScreenshotOptions(browser.ScreenshotWithFormat(format))
	if format == "jpeg" {
		opts = browser.NewScreenshotOptions(
			browser.ScreenshotWithFormat(format),
			browser.ScreenshotWithQuality(quality),
		)
	}
	data, err := s.browser.Screenshot(tabID, opts)
	if err != nil {
		return toolError(err)
	}
	return toolImage(fmt.Sprintf("image/%s", format), data)
}

func (s *BrowserMcpServer) handleGetTexts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	texts, err := s.browser.GetTexts(loc.TabID, loc.Target, loc.TargetType)
	if err != nil {
		return toolError(err)
	}
	return toolJSON(texts)
}

func (s *BrowserMcpServer) handleGetHtmls(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	htmls, err := s.browser.GetHtmls(loc.TabID, loc.Target, loc.TargetType)
	if err != nil {
		return toolError(err)
	}
	return toolJSON(htmls)
}
