package param

import "github.com/mark3labs/mcp-go/mcp"

// Args extracts tool call arguments as a map. Nil arguments yield an empty map.
func Args(request mcp.CallToolRequest) map[string]any {
	args := request.GetArguments()
	if args == nil {
		return map[string]any{}
	}
	return args
}
