## Context

`libs/browser` 已通过 `Browser` 接口与 `BrowserApiHandler` / `TabHandler` 暴露 REST（Tab CRUD、导航、元素操作、snapshot/screenshot、cookies、Remote WS）。Agent 侧主流接入是 MCP；本仓库尚无 MCP 实现。opsys server 已有 mark3labs 用法（SSE + Streamable HTTP）可作挂载参考，但参数规范与工具面不能直接依赖 opsys。

已确认约束：

| 项 | 决策 |
|---|---|
| 工具名 | Playwright 风格 `browser_*` |
| 参数名 | Playwright 风格 **camelCase**（`tabId`、`targetType`、`doubleClick`） |
| Tab 上下文 | 几乎所有 tab 操作**必填** `tabId` |
| 定位 | `target` + `targetType`；snapshot 后优先引导 `ref` |
| 范围 | P0 + P1；不做 wait/evaluate/dialog/network/console；不做 Remote/Screencast |
| 传输 | HTTP/SSE（挂 browserd）；同时挂 Streamable HTTP（与 opsys 一致，便于新 Client） |
| 参数框架 | `libs/browser/mcp/param` 轻量解析库 |

## Goals / Non-Goals

**Goals:**

- 以 `Browser` 接口为唯一执行后端，组装 mark3labs `MCPServer` 与工具 handler
- 提供标准化参数解析（必填/可选、enum、默认值、类型加工），错误信息可直接返回给 Client
- 在 `browserd` 上挂载 MCP HTTP 端点（SSE + Streamable），与现有 REST 共存
- 工具面覆盖 P0（现有 REST）+ P1（接口已有：switch tab、upload files、press key、print PDF）
- snapshot 工具 description / 返回引导 Agent 使用 `ref` + `targetType=ref`

**Non-Goals:**

- stdio transport（本期不做；框架不阻塞后续加）
- wait_for / evaluate / dialog / network / console / fill_form / drag
- Remote WebSocket / screencast 作为 MCP 工具
- 修改现有 REST JSON 字段名或行为
- 鉴权体系重构（复用 browserd 已有 JWT 中间件即可，若已配置）
- 依赖或复制 opsys `resource/mcp/param` 包

## Decisions

### 1. 包结构与依赖方向

**选择**:

```
libs/browser/mcp/
  server.go          # NewBrowserMcpServer / RegisterTools / HTTP 挂载 helpers
  tools_*.go         # 按域拆分工具注册与 handler
  param/
    string.go        # RequiredString / OptionalString / StringEnum
    number.go        # OptionalInt / OptionalFloat / Bool
    args.go          # 从 mcp.CallToolRequest 取 Arguments
    errors.go        # 统一错误文案
```

- MCP 层只依赖 `browser.Browser`（及本包 param），**不**经 HTTP 自调用 REST
- `browserd` 在现有 router 旁挂载 MCP handler

**备选**: 独立顶层 `libs/mcpparam` —— 过早抽象，否决；先落在 `libs/browser/mcp/param`，日后可上提。

**备选**: MCP handler 内部 HTTP client 调 REST —— 多一层序列化与鉴权复杂度，否决。

### 2. Transport：HTTP/SSE + Streamable HTTP

**选择**:

| 路径（相对 MCP 前缀，默认 `/mcp`） | 实现 |
|---|---|
| `/mcp/sse` + `/mcp/sse/message` | mark3labs `SSEServer` |
| `/mcp` 与 `/mcp/streamable` | mark3labs `StreamableHTTPServer` |

前缀可由环境变量 `MCP_PATH`（默认 `/mcp`）配置。挂载方式对齐 opsys `apimodel.McpServer` 的双端点模式，但实现保持在 goutils 内、更薄。

**备选**: 仅 SSE —— 部分新 Client 偏好 Streamable，成本低故两者都挂。  
**备选**: 仅 stdio —— 与「挂 browserd」决策不符，否决。

### 3. 工具命名与拆分策略

**选择**: Playwright 风格前缀 `browser_`，但 **Tab 管理拆成独立工具**（而非单一 `browser_tabs` + `action`），以便强制 `tabId` 与清晰 schema：

| 工具 | 映射 | 必填关键参数 |
|---|---|---|
| `browser_list_tabs` | `ListTabs` | — |
| `browser_new_tab` | `NewTab` | 可选 `url`/`width`/`height` |
| `browser_close_tab` | `CloseTab` | `tabId` |
| `browser_get_tab` | `GetTab` | `tabId` |
| `browser_select_tab` | `SwitchToTab` | `tabId` |
| `browser_navigate` | `Navigate` | `tabId`, `url` |
| `browser_navigate_back` | `GoBack` | `tabId` |
| `browser_navigate_forward` | `GoForward` | `tabId` |
| `browser_reload` | `Reload` | `tabId` |
| `browser_click` | `Click` | `tabId`, `target`, `targetType`；可选 `button`, `doubleClick`/`count` |
| `browser_type` | `Type` | `tabId`, `target`, `targetType`, `text` |
| `browser_hover` | `Hover` | `tabId`, `target`, `targetType` |
| `browser_select_option` | `SelectOption` | `tabId`, `target`, `targetType`, `values`；可选 `optionType`, `selected` |
| `browser_snapshot` | `Snapshot` | `tabId`；可选 `type`（`a11y`/`dom`，默认 `a11y`） |
| `browser_take_screenshot` | `Screenshot` | `tabId`；可选 format/quality/target 等 |
| `browser_get_texts` | `GetTexts` | `tabId`, `target`, `targetType` |
| `browser_get_htmls` | `GetHtmls` | `tabId`, `target`, `targetType` |
| `browser_cookie_list` | `GetCookies` | `tabId` |
| `browser_cookie_set` | `SetCookies`（单 cookie） | `tabId`, cookie 字段 |
| `browser_file_upload` | `SetInputFiles` | `tabId`, `target`, `targetType`, `paths` |
| `browser_press_key` | `DispatchKeyEvent`（合成 press=down+up） | `tabId`, `key`；可选 `modifiers` |
| `browser_pdf_save` | `PrintToPdf` | `tabId` |

与 Playwright 差异（有意保留，已确认）：

- 强制 `tabId`（Playwright 用隐式 current page）
- 保留 `targetType`（Playwright 常把 ref/selector 塞进单一 `target`）
- 有 `browser_navigate_forward` / `browser_reload` / `browser_get_texts` / `browser_get_htmls`（我们已有能力）
- 无 `browser_wait_for` / evaluate / network 等（Non-Goals）

### 4. 参数标准化（camelCase + param 库）

**选择**: 工具 schema 与解析层统一 **camelCase**，与 Playwright / 用户决策 2C 一致。

标准参数约定：

| 语义 | 参数名 | MCP 类型 | 规则 |
|---|---|---|---|
| Tab ID | `tabId` | string | 非空；无效 tab → tool error |
| 定位目标 | `target` | string | 与 `targetType` 成对 |
| 定位类型 | `targetType` | string | enum: `ref`, `css-selector`, `xpath`；缺省时若像 snapshot ref 可文档推荐显式传 `ref`，**解析层不猜**（避免歧义） |
| URL | `url` | string | 非空 |
| 双击 | `doubleClick` | boolean | true 时映射 `count=2`；与 `count` 同时出现时 `count` 优先 |
| 点击次数 | `count` | number | 默认 1 |
| 鼠标按钮 | `button` | string | enum 对齐 `MouseButton*`；默认 `left` |
| 下拉选项 | `values` | string[] | 映射 `SelectOption` 的 options；`optionType` 默认 `text` |
| 截图格式 | `type`（screenshot） / `format` | string | 与现有 `ScreenOptionsQuery` 对齐；默认 `png`（对齐 Playwright screenshot） |
| Snapshot 类型 | `type` | string | `a11y` \| `dom`；默认 `a11y` |
| 按键 | `key` | string | press 时服务端 down+up |
| 上传路径 | `paths` | string[] | 非空 |

`param` 包提供：

- `Args(request)` → `map[string]any`
- `RequiredString` / `OptionalString` / `StringEnum`
- `OptionalBool` / `OptionalInt` / `RequiredStringSlice`
- 错误格式：`{param} is required` / `{param} must be one of: ...`

**备选**: snake_case（opsys）—— 与决策 2C 冲突，否决。  
**备选**: 解析层自动推断 targetType —— 歧义风险高，否决；用 description 引导。

### 5. 结果形态

**选择**:

- 结构化数据（tabs、cookie、texts/htmls）：`mcp.NewToolResultText` + JSON（indent 可选，默认 compact JSON）
- snapshot：纯文本（a11y/dom 字符串）；description 注明「用返回的 ref 调用动作工具，且 `targetType=ref`」
- screenshot / pdf：`mcp.NewToolResultImage` / 或 base64 text resource；优先 mark3labs 支持的 image/blob 结果类型，若版本不便则 base64 + mime 文本
- 无数据成功（navigate/click 等）：短文本 `ok` 或含 `tabId` 的 JSON `{ "ok": true, "tabId": "..." }`

错误：一律 `mcp.NewToolResultError(...)`，不把业务错误升成 protocol-level error（除非参数解析前 panic 级故障）。

### 6. browserd 集成

**选择**:

```go
mcpHandler := browsermcp.NewHTTPHandler(browser, browsermcp.WithPath(mcpPath))
// 与 REST router 组合：StripPrefix 或 chi/mux 挂载
```

- `instance` 模式：一个 `Browser` → 一个 MCP server
- `pool` 模式：首期仅 **instance** 挂 MCP；pool 若需 MCP，后续按 pool API 设计（Open Question 可记「暂不支持 pool MCP」）
- 若 `JWT_SECRET` 已启用，MCP 路径走同一 JWT 中间件（与 REST 一致）

### 7. 测试策略

- `param`：纯单测（边界、enum、默认值）
- tools：`Browser` stub（复用/扩展 `remote_test.go` 的 stub 模式）验证参数映射与错误路径
- 不强制 e2e 启真实 Chrome（可选用现有 rod 集成测试标签）

## Risks / Trade-offs

- [与 Playwright 工具/参数不完全一致] → Mitigation：在工具 description 写明差异（必填 `tabId`、`targetType`）；不追求 100% 兼容
- [camelCase 与 opsys snake_case 分裂] → Mitigation：browser MCP 独立规范；文档标明；param 库不跨仓共享
- [screenshot/pdf 大 payload 撑爆 context] → Mitigation：默认合理 quality；description 建议优先 snapshot；后续可加 max 尺寸
- [pool 模式未挂 MCP] → Mitigation：文档标明首期 instance；需要时另开 change
- [press_key 语义弱于 Playwright] → Mitigation：仅封装 down+up；复杂组合键用 `modifiers`；不做 type_text 全页输入

## Migration Plan

1. 添加 `mark3labs/mcp-go` 依赖
2. 落地 `param` → tools → server → browserd 挂载
3. 本地用 MCP Inspector / Cursor 连 `http://localhost:8080/mcp/sse` 验证
4. 回滚：去掉 MCP 挂载与依赖即可，REST 不受影响

## Open Questions

- pool 模式是否在本期用「默认 browser」挂一套 MCP？（默认：**不做**，仅 instance）
- screenshot 默认 `png`（对齐 Playwright）还是 `jpeg`（对齐现有 REST query 默认）？**建议 MCP 默认 png**，REST 不变
- `browser_cookie_set` 字段是扁平（`name`/`value`/`domain`…）还是嵌套 `cookie` 对象？**建议扁平**，贴近 Playwright storage 工具
