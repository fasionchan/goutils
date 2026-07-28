## Why

`libs/browser` 已提供完整的浏览器控制能力（Tab 管理、导航、元素操作、snapshot/screenshot、cookies 等）及 HTTP REST API，但 Agent 生态主流接入方式是 MCP。需要在现有 `Browser` 接口之上暴露 MCP 工具面，使 Cursor / 其他 MCP Client 可直接驱动 browserd，而无需手写 REST 调用。

## What Changes

- 新增基于 **mark3labs/mcp-go** 的 Browser MCP Server（包路径 `libs/browser/mcp`）
- 工具命名对齐 Playwright MCP：`browser_*`；参数命名对齐 Playwright：**camelCase**（如 `tabId`、`targetType`、`doubleClick`）
- 几乎所有 tab 作用域工具**必填** `tabId`（多 tab / pool 场景明确）
- 定位模型保持现有 `target` + `targetType`（`ref` | `css-selector` | `xpath`）；snapshot 后优先引导 `ref` 工作流
- 首期工具范围 **P0 + P1**（覆盖现有 REST + `Browser` 接口已有但 REST 未全暴露的能力）；不做 wait/evaluate/dialog/network/console；不做 Remote/Screencast
- 传输：**HTTP/SSE**（及 mark3labs 的 Streamable HTTP，若框架支持），挂载到 `browserd`
- 新增轻量参数解析框架（`libs/browser/mcp/param`）：类型校验、默认值、enum、加工处理；风格参考 opsys `resource/mcp/param`，但不依赖 opsys
- 新增依赖：`github.com/mark3labs/mcp-go`

## Capabilities

### New Capabilities

- `browser-mcp`: Browser MCP 工具清单、服务组装、HTTP/SSE 挂载到 browserd，以及与 `Browser` 接口的调用映射
- `browser-mcp-param`: MCP 工具参数标准化解析库（camelCase 命名、类型、默认值、enum、错误消息）

### Modified Capabilities

（无）

## Impact

- **代码**: 新增 `libs/browser/mcp/`（server、tools、param）；`libs/browser/cmd/browserd/main.go` 挂载 MCP HTTP 路由；必要时小范围补齐 `Browser`/`TabHandler` 已有能力的薄封装（不改无关 REST 契约）
- **依赖**: 新增 `github.com/mark3labs/mcp-go`
- **运行时**: browserd 额外暴露 MCP SSE（及可选 Streamable HTTP）端点；现有 `/Tabs` REST 与 Remote WebSocket 保持不变
- **兼容性**: 纯新增能力，无 **BREAKING** 变更
