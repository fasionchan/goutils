# browser-mcp

基于 mark3labs/mcp-go 的 Browser MCP 服务：工具面、HTTP SSE/Streamable 挂载，以及与 Browser 接口的调用映射。

## Requirements

### Requirement: 基于 mark3labs 组装 Browser MCP Server

系统 MUST 在 `libs/browser/mcp` 使用 `github.com/mark3labs/mcp-go` 组装 MCP Server，并以 `browser.Browser` 为执行后端（MUST NOT 通过 HTTP 自调用 REST）。

#### Scenario: 成功创建 server
- **WHEN** 调用方传入非 nil 的 `Browser` 创建 MCP Server
- **THEN** 返回已注册 P0+P1 工具的可用 server

#### Scenario: Browser 为 nil
- **WHEN** 以 nil `Browser` 创建 server
- **THEN** 创建失败或在首次 tool 调用时返回明确错误（实现二选一，但 MUST 可观测）

### Requirement: HTTP SSE 与 Streamable 挂载

系统 MUST 将 MCP Server 以 HTTP 方式挂载，至少提供 SSE 端点与 Streamable HTTP 端点；默认路径前缀为 `/mcp`，可配置。

#### Scenario: browserd instance 模式挂载
- **WHEN** 以 instance 模式启动 browserd 且启用 MCP
- **THEN** Client 可通过 `/mcp/sse`（及对应 message 路径）和 `/mcp`（或 `/mcp/streamable`）建立 MCP 会话

#### Scenario: 与 REST 共存
- **WHEN** browserd 已挂载 `/Tabs` REST 与 MCP
- **THEN** 两类端点同时可用，互不替换

### Requirement: Tab 管理工具

系统 MUST 提供以下工具：`browser_list_tabs`、`browser_new_tab`、`browser_close_tab`、`browser_get_tab`、`browser_select_tab`。

#### Scenario: 列出 tabs
- **WHEN** 调用 `browser_list_tabs`
- **THEN** 返回当前 tabs 的结构化结果（至少含 id/title/url）

#### Scenario: 新建 tab
- **WHEN** 调用 `browser_new_tab` 并可选提供 `url`
- **THEN** 创建 tab 并返回含新 `tabId` 的结果

#### Scenario: 关闭 tab 缺少 tabId
- **WHEN** 调用 `browser_close_tab` 且未提供 `tabId`
- **THEN** 返回 tool error，指出 `tabId is required`

#### Scenario: 切换 tab
- **WHEN** 调用 `browser_select_tab` 并提供有效 `tabId`
- **THEN** 调用 `Browser.SwitchToTab` 成功并返回 ok 类结果

### Requirement: 导航工具强制 tabId

系统 MUST 提供 `browser_navigate`、`browser_navigate_back`、`browser_navigate_forward`、`browser_reload`；除 list/new 外，这些工具 MUST 要求 `tabId`。

#### Scenario: 导航
- **WHEN** 调用 `browser_navigate` 并提供 `tabId` 与 `url`
- **THEN** 目标 tab 导航到该 URL

#### Scenario: 前进/后退/刷新
- **WHEN** 分别调用 `browser_navigate_back` / `browser_navigate_forward` / `browser_reload` 并提供有效 `tabId`
- **THEN** 对应 `GoBack` / `GoForward` / `Reload` 被执行

### Requirement: 元素动作工具使用 target 与 targetType

系统 MUST 提供 `browser_click`、`browser_type`、`browser_hover`、`browser_select_option`、`browser_file_upload`；这些工具 MUST 要求 `tabId`、`target`、`targetType`（`ref` | `css-selector` | `xpath`）。

#### Scenario: 以 ref 点击
- **WHEN** 调用 `browser_click`，`targetType` 为 `ref`，`target` 为 snapshot ref
- **THEN** 后端以 ref 定位并执行点击

#### Scenario: doubleClick 映射
- **WHEN** 调用 `browser_click` 且 `doubleClick` 为 true、未提供 `count`
- **THEN** 以点击次数 2 调用 `Browser.Click`

#### Scenario: 输入文本
- **WHEN** 调用 `browser_type` 并提供 `text`
- **THEN** 在目标元素上执行 `Type`

#### Scenario: 上传文件
- **WHEN** 调用 `browser_file_upload` 并提供非空 `paths`
- **THEN** 调用 `SetInputFiles` 将文件设置到目标元素

#### Scenario: 下拉选择
- **WHEN** 调用 `browser_select_option` 并提供 `values`
- **THEN** 调用 `SelectOption`；未提供 `optionType` 时默认按 `text` 处理

### Requirement: 观察工具 snapshot 与 screenshot

系统 MUST 提供 `browser_snapshot` 与 `browser_take_screenshot`，均要求 `tabId`。

#### Scenario: 默认 a11y snapshot
- **WHEN** 调用 `browser_snapshot` 且未指定 `type`
- **THEN** 返回 a11y snapshot 文本，且工具描述或返回说明引导使用 `targetType=ref` 进行后续动作

#### Scenario: 截图
- **WHEN** 调用 `browser_take_screenshot` 并提供有效 `tabId`
- **THEN** 返回图像类 MCP 结果（或带 mime 的 base64），默认格式为 png（除非调用方指定其他格式）

### Requirement: 文本与 HTML 提取工具

系统 MUST 提供 `browser_get_texts` 与 `browser_get_htmls`，要求 `tabId`、`target`、`targetType`。

#### Scenario: 获取 texts
- **WHEN** 调用 `browser_get_texts` 定位到匹配元素
- **THEN** 返回文本列表的结构化结果

### Requirement: Cookie 工具

系统 MUST 提供 `browser_cookie_list` 与 `browser_cookie_set`，要求 `tabId`。

#### Scenario: 列出 cookies
- **WHEN** 调用 `browser_cookie_list`
- **THEN** 返回该 tab 的 cookie 列表

#### Scenario: 设置 cookie
- **WHEN** 调用 `browser_cookie_set` 并提供至少 `name` 与 `value`
- **THEN** 在目标 tab 设置 cookie 并返回结果

### Requirement: 按键与 PDF 工具

系统 MUST 提供 `browser_press_key`（基于 `DispatchKeyEvent` 合成 down+up）与 `browser_pdf_save`（基于 `PrintToPdf`），均要求 `tabId`。

#### Scenario: 按键
- **WHEN** 调用 `browser_press_key` 并提供 `key`
- **THEN** 向目标 tab 派发键盘 down 与 up

#### Scenario: 导出 PDF
- **WHEN** 调用 `browser_pdf_save`
- **THEN** 返回该 tab 的 PDF 字节结果（适合 MCP 传输的编码形式）

### Requirement: 工具错误不破坏会话

业务失败（未知 tab、定位失败、后端错误）MUST 以 MCP tool error 结果返回，MUST NOT 无故断开 SSE/Streamable 会话。

#### Scenario: 未知 tabId
- **WHEN** 任一要求 `tabId` 的工具收到不存在的 id
- **THEN** 返回 tool error，会话保持可用

### Requirement: 不暴露 Remote 与非目标能力

MCP 工具面 MUST NOT 包含 screencast/remote WebSocket 控制，也 MUST NOT 在本期包含 wait_for、evaluate、dialog、network、console 类工具。

#### Scenario: 工具列表审计
- **WHEN** Client 列出 tools
- **THEN** 列表包含本 spec 所列 P0+P1 工具，且不包含 remote/screencast/wait_for/evaluate/dialog/network/console 工具名
