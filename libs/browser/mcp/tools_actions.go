package browsermcp

import (
	"context"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func locatorToolOptions(extra ...mcp.ToolOption) []mcp.ToolOption {
	base := []mcp.ToolOption{
		mcp.WithString("tabId", mcp.Required(), mcp.Description("Tab ID")),
		mcp.WithString("target", mcp.Required(), mcp.Description("Element locator value (prefer snapshot ref)")),
		mcp.WithString("targetType", mcp.Required(), mcp.Description("Locator type"), mcp.Enum(targetTypeAllowed...)),
	}
	return append(base, extra...)
}

func (s *BrowserMcpServer) registerActionTools() {
	clickOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Click an element. Prefer targetType=ref from browser_snapshot. Requires tabId."),
	}, locatorToolOptions(
		mcp.WithString("button", mcp.Description("Mouse button"), mcp.Enum(buttonAllowed...), mcp.DefaultString(browser.MouseButtonLeft)),
		mcp.WithBoolean("doubleClick", mcp.Description("Double click when true and count omitted"), mcp.DefaultBool(false)),
		mcp.WithNumber("count", mcp.Description("Click count; overrides doubleClick when set"), mcp.DefaultNumber(1)),
	)...)
	s.AddTool(mcp.NewTool("browser_click", clickOpts...), s.handleClick)

	typeOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Type text into an element. Prefer targetType=ref from browser_snapshot. Requires tabId."),
	}, locatorToolOptions(
		mcp.WithString("text", mcp.Required(), mcp.Description("Text to type")),
	)...)
	s.AddTool(mcp.NewTool("browser_type", typeOpts...), s.handleType)

	hoverOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Hover over an element. Prefer targetType=ref from browser_snapshot. Requires tabId."),
	}, locatorToolOptions()...)
	s.AddTool(mcp.NewTool("browser_hover", hoverOpts...), s.handleHover)

	selectOpts := append([]mcp.ToolOption{
		mcp.WithDescription("Select option(s) in a dropdown. Prefer targetType=ref. Requires tabId."),
	}, locatorToolOptions(
		mcp.WithArray("values", mcp.Required(), mcp.Description("Option values to select"), mcp.WithStringItems()),
		mcp.WithString("optionType", mcp.Description("How to match options"), mcp.Enum(optionTypeAllowed...), mcp.DefaultString(browser.OptionLocatorTypeText)),
		mcp.WithBoolean("selected", mcp.Description("Select (true) or deselect (false)"), mcp.DefaultBool(true)),
	)...)
	s.AddTool(mcp.NewTool("browser_select_option", selectOpts...), s.handleSelectOption)
}

func (s *BrowserMcpServer) handleClick(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	button, err := param.OptionalStringEnum(args, "button", buttonAllowed, browser.MouseButtonLeft)
	if err != nil {
		return toolError(err)
	}
	count := 1
	if param.Has(args, "count") {
		count, err = param.OptionalInt(args, "count", 1)
		if err != nil {
			return toolError(err)
		}
	} else {
		doubleClick, err := param.OptionalBool(args, "doubleClick", false)
		if err != nil {
			return toolError(err)
		}
		if doubleClick {
			count = 2
		}
	}
	if err := s.browser.Click(loc.TabID, loc.Target, loc.TargetType, button, count); err != nil {
		return toolError(err)
	}
	return toolOK(loc.TabID)
}

func (s *BrowserMcpServer) handleType(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	text, err := param.RequiredString(args, "text")
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.Type(loc.TabID, loc.Target, loc.TargetType, text); err != nil {
		return toolError(err)
	}
	return toolOK(loc.TabID)
}

func (s *BrowserMcpServer) handleHover(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.Hover(loc.TabID, loc.Target, loc.TargetType); err != nil {
		return toolError(err)
	}
	return toolOK(loc.TabID)
}

func (s *BrowserMcpServer) handleSelectOption(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := param.Args(request)
	loc, err := parseLocator(args)
	if err != nil {
		return toolError(err)
	}
	values, err := param.RequiredStringSlice(args, "values")
	if err != nil {
		return toolError(err)
	}
	optionType, err := param.OptionalStringEnum(args, "optionType", optionTypeAllowed, browser.OptionLocatorTypeText)
	if err != nil {
		return toolError(err)
	}
	selected, err := param.OptionalBool(args, "selected", true)
	if err != nil {
		return toolError(err)
	}
	if err := s.browser.SelectOption(loc.TabID, loc.Target, loc.TargetType, values, optionType, selected); err != nil {
		return toolError(err)
	}
	return toolOK(loc.TabID)
}
