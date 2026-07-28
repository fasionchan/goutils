# browser-mcp-param

Browser MCP 工具参数标准化解析库（camelCase 命名、类型、默认值、enum、错误消息）。

## Requirements

### Requirement: MCP 参数命名使用 camelCase

`libs/browser/mcp/param` 解析的参数名 MUST 使用 camelCase（例如 `tabId`、`targetType`、`doubleClick`），与 Browser MCP 工具 schema 一致。

#### Scenario: 按标准名读取参数
- **WHEN** handler 通过 param 库读取 `tabId`
- **THEN** 库从 `CallToolRequest` 的 arguments 中按键名 `tabId` 取值，且 MUST NOT 要求调用方使用 snake_case 别名

### Requirement: 必填字符串与可选字符串

param 库 MUST 提供必填/可选字符串读取能力；必填缺失或为空白时 MUST 返回可直接展示给 Client 的错误。

#### Scenario: 缺少必填 tabId
- **WHEN** 调用 `RequiredString(args, "tabId")` 且 arguments 无 `tabId` 或为空字符串
- **THEN** 返回错误，消息包含 `tabId is required`

#### Scenario: 可选字符串缺省
- **WHEN** 调用 `OptionalString(args, "url")` 且未提供 `url`
- **THEN** 返回空字符串与 nil error

### Requirement: 枚举字符串校验

param 库 MUST 支持字符串 enum 校验；非法值 MUST 列出允许集合。

#### Scenario: targetType 非法
- **WHEN** 以允许集 `ref,css-selector,xpath` 解析 `targetType` 且值为 `foo`
- **THEN** 返回错误，消息包含 `targetType` 与允许值列表

#### Scenario: targetType 合法
- **WHEN** `targetType` 为 `ref`
- **THEN** 返回 `"ref"` 与 nil error

### Requirement: 布尔与整数可选参数

param 库 MUST 支持从 JSON number/bool（及常见可兼容形态）解析可选 bool/int，并支持默认值。

#### Scenario: doubleClick 缺省
- **WHEN** 未提供 `doubleClick` 且默认值为 false
- **THEN** 解析结果为 false，无 error

#### Scenario: count 由 number 提供
- **WHEN** arguments 含 `"count": 2`
- **THEN** `OptionalInt` 返回 `2`

### Requirement: 字符串切片参数

param 库 MUST 支持读取 JSON 字符串数组（如 `paths`、`values`）；必填切片为空时 MUST 报错。

#### Scenario: paths 必填非空
- **WHEN** `RequiredStringSlice(args, "paths")` 收到 `[]` 或缺失
- **THEN** 返回错误，消息包含 `paths`

#### Scenario: values 合法数组
- **WHEN** arguments 含 `"values": ["a","b"]`
- **THEN** 返回长度为 2 的字符串切片

### Requirement: 统一从 CallToolRequest 提取 Arguments

param 库 MUST 提供从 mark3labs `mcp.CallToolRequest` 提取 `map[string]any`（或等价结构）的入口，供上述解析函数使用。

#### Scenario: 空 arguments
- **WHEN** request 无 arguments
- **THEN** 提取结果为空 map，后续 Required* 按缺失处理
