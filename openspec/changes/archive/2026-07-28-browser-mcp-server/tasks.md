## 1. 依赖与包骨架

- [x] 1.1 在 go.mod 添加 `github.com/mark3labs/mcp-go` 并 tidy
- [x] 1.2 创建 `libs/browser/mcp/` 与 `libs/browser/mcp/param/` 包骨架（空文件/最小类型即可）

## 2. 参数解析库（browser-mcp-param）

- [x] 2.1 实现从 `mcp.CallToolRequest` 提取 arguments 的入口
- [x] 2.2 实现 `RequiredString` / `OptionalString` / `StringEnum`（camelCase 键名）
- [x] 2.3 实现 `OptionalBool` / `OptionalInt`（含默认值）与 `RequiredStringSlice`
- [x] 2.4 为 param 包补充单测（缺失、enum、空切片、默认值）

## 3. MCP Server 与工具注册框架

- [x] 3.1 实现 `NewBrowserMcpServer(browser Browser, opts...)`：创建 mark3labs server 并预留 RegisterTools
- [x] 3.2 定义共享 helper：定位参数解析（`tabId`+`target`+`targetType`）、统一 ok/JSON/error 结果构造
- [x] 3.3 实现 HTTP 挂载 helper：SSE（`/sse` + message）与 Streamable（`/` 与 `/streamable`），支持可配置前缀（默认 `/mcp`）

## 4. P0 工具实现

- [x] 4.1 注册 Tab 管理：`browser_list_tabs`、`browser_new_tab`、`browser_close_tab`、`browser_get_tab`
- [x] 4.2 注册导航：`browser_navigate`、`browser_navigate_back`、`browser_navigate_forward`、`browser_reload`
- [x] 4.3 注册动作：`browser_click`（含 `doubleClick`/`count`/`button`）、`browser_type`、`browser_hover`、`browser_select_option`
- [x] 4.4 注册观察：`browser_snapshot`（默认 a11y + ref 引导）、`browser_take_screenshot`（默认 png）、`browser_get_texts`、`browser_get_htmls`
- [x] 4.5 注册 Cookie：`browser_cookie_list`、`browser_cookie_set`（扁平字段）

## 5. P1 工具实现

- [x] 5.1 注册 `browser_select_tab` → `SwitchToTab`
- [x] 5.2 注册 `browser_file_upload` → `SetInputFiles`（`paths`）
- [x] 5.3 注册 `browser_press_key` → `DispatchKeyEvent`（down+up，可选 `modifiers`）
- [x] 5.4 注册 `browser_pdf_save` → `PrintToPdf`（合适的 MCP 二进制/base64 结果）

## 6. browserd 集成与测试

- [x] 6.1 在 `libs/browser/cmd/browserd/main.go` instance 模式挂载 MCP HTTP；复用现有 JWT 中间件（若配置）
- [x] 6.2 使用 Browser stub 为关键工具编写 handler 单测（缺 `tabId`、映射、错误路径）
- [x] 6.3 工具列表审计测试：确保含 P0+P1，且不含 remote/wait_for/evaluate/dialog/network/console
- [x] 6.4 手动或脚本验证：启动 browserd 后 MCP Inspector/Cursor 可连 SSE 并完成 list_tabs → new_tab → navigate → snapshot → click
