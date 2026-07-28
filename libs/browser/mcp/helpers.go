package browsermcp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fasionchan/goutils/libs/browser"
	"github.com/fasionchan/goutils/libs/browser/mcp/param"
	"github.com/mark3labs/mcp-go/mcp"
)

var targetTypeAllowed = []string{
	browser.LocatorTypeRef,
	browser.LocatorTypeCssSelector,
	browser.LocatorTypeXPath,
}

var buttonAllowed = []string{
	browser.MouseButtonLeft,
	browser.MouseButtonMiddle,
	browser.MouseButtonRight,
	browser.MouseButtonBack,
	browser.MouseButtonForward,
	browser.MouseButtonNone,
}

var snapshotTypeAllowed = []string{
	browser.SnapshotTypeA11y,
	browser.SnapshotTypeDom,
}

var optionTypeAllowed = []string{
	browser.OptionLocatorTypeText,
	browser.OptionLocatorTypeCssSelector,
	browser.OptionLocatorTypeRegex,
}

type locatorParams struct {
	TabID      string
	Target     string
	TargetType string
}

func parseTabID(args map[string]any) (string, error) {
	return param.RequiredString(args, "tabId")
}

func parseLocator(args map[string]any) (*locatorParams, error) {
	tabID, err := parseTabID(args)
	if err != nil {
		return nil, err
	}
	target, err := param.RequiredString(args, "target")
	if err != nil {
		return nil, err
	}
	targetType, err := param.StringEnum(args, "targetType", targetTypeAllowed)
	if err != nil {
		return nil, err
	}
	return &locatorParams{TabID: tabID, Target: target, TargetType: targetType}, nil
}

func toolError(err error) (*mcp.CallToolResult, error) {
	if err == nil {
		return mcp.NewToolResultError("unknown error"), nil
	}
	return mcp.NewToolResultError(err.Error()), nil
}

func toolOK(tabID string) (*mcp.CallToolResult, error) {
	payload := map[string]any{"ok": true}
	if tabID != "" {
		payload["tabId"] = tabID
	}
	return toolJSON(payload)
}

func toolJSON(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal result failed: %v", err)), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func toolText(text string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(text), nil
}

func toolImage(mimeType string, data []byte) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultImage(mimeType, base64.StdEncoding.EncodeToString(data), mimeType), nil
}

func toolPDF(data []byte) (*mcp.CallToolResult, error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	return mcp.NewToolResultResource("application/pdf", mcp.BlobResourceContents{
		URI:      "browser://pdf",
		MIMEType: "application/pdf",
		Blob:     encoded,
	}), nil
}

func cookieFromArgs(args map[string]any) (*http.Cookie, error) {
	name, err := param.RequiredString(args, "name")
	if err != nil {
		return nil, err
	}
	value, err := param.RequiredString(args, "value")
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{Name: name, Value: value}
	if domain, err := param.OptionalString(args, "domain"); err != nil {
		return nil, err
	} else {
		c.Domain = domain
	}
	if path, err := param.OptionalString(args, "path"); err != nil {
		return nil, err
	} else if path != "" {
		c.Path = path
	}
	secure, err := param.OptionalBool(args, "secure", false)
	if err != nil {
		return nil, err
	}
	c.Secure = secure
	httpOnly, err := param.OptionalBool(args, "httpOnly", false)
	if err != nil {
		return nil, err
	}
	c.HttpOnly = httpOnly
	return c, nil
}

func ensurePathPrefix(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return DefaultMCPPath
	}
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
