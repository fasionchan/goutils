## MODIFIED Requirements

### Requirement: 按字节上限截断 Markdown 内容

`TruncateMarkdown` SHALL 接受 Markdown 字符串与 `maxBytes` 上限，返回截断后的字符串及是否发生截断的布尔值。字节计数 MUST 按 UTF-8 编码计算，与 Go `len([]byte(s))` 一致。`TruncateMarkdown` MUST NOT 追加任何后缀；`maxBytes` 全部用于内容截断预算。

#### Scenario: 内容未超出上限

- **WHEN** 输入 Markdown 的 UTF-8 字节长度小于或等于 `maxBytes`
- **THEN** 返回原始内容不变，截断标志为 `false`

#### Scenario: 内容超出上限

- **WHEN** 输入 Markdown 的 UTF-8 字节长度大于 `maxBytes`
- **THEN** 返回长度不超过 `maxBytes` 的 Markdown 字符串，截断标志为 `true`，且不包含任何后缀

#### Scenario: maxBytes 非正数

- **WHEN** `maxBytes` 小于或等于 0
- **THEN** 返回空字符串，截断标志为 `true`（若原内容非空）或 `false`（若原内容为空）

### Requirement: 在块级节点边界截断

截断 MUST 在 goldmark AST 块级节点边界进行，不得将单个块级节点从中间切开。若某块级节点本身超过 `maxBytes`，该节点 MUST 作为整体保留或整体丢弃，不得产生不完整的 Markdown 结构。

#### Scenario: 多段落截断

- **WHEN** Markdown 包含多个段落且总长度超出 `maxBytes`
- **THEN** 结果包含从文档开头起、完整保留的连续块级节点，不包含被截断节点的不完整片段

#### Scenario: 代码块不截断中间

- **WHEN** Markdown 包含 fenced code block 且截断点落在该代码 block 内
- **THEN** 要么完整保留该代码块（若剩余空间足够），要么不包含该代码块，不得输出未闭合的 fence 标记

#### Scenario: 列表项保持完整

- **WHEN** Markdown 包含有序或无序列表且截断点落在某列表项内
- **THEN** 该列表项 MUST 作为整体保留或整体丢弃

### Requirement: 截断时追加省略标记

`TruncateMarkdownWithSuffix` SHALL 在 `TruncateMarkdown` 基础上追加省略标记。`TruncateMarkdownWithSuffix` MUST 调用 `TruncateMarkdown` 完成内容截断，而非自行实现截断逻辑。当 `suffix` 为空字符串时，MUST 使用默认省略标记 `...`。后缀字节长度 MUST 计入 `maxBytes` 总预算（先从预算中扣除后缀长度，再将剩余预算传给 `TruncateMarkdown`）。

#### Scenario: 默认省略标记

- **WHEN** 通过 `TruncateMarkdownWithSuffix` 截断且 `suffix` 为空字符串
- **THEN** 结果以 `...` 结尾，且总字节长度不超过 `maxBytes`

#### Scenario: 自定义省略标记

- **WHEN** 通过 `TruncateMarkdownWithSuffix` 截断且指定非空 `suffix`
- **THEN** 结果以指定 `suffix` 结尾，且总字节长度不超过 `maxBytes`

#### Scenario: 省略标记占用预算

- **WHEN** `maxBytes` 较小，扣除省略标记后无法容纳任何块级节点
- **THEN** 返回仅含省略标记（或空字符串，若 `maxBytes` 不足以容纳省略标记）的结果，截断标志为 `true`

#### Scenario: TruncateMarkdown 不追加后缀

- **WHEN** 直接调用 `TruncateMarkdown` 且内容被截断
- **THEN** 返回结果不包含任何省略标记

### Requirement: AST 解析失败时的降级策略

当 goldmark AST 解析或遍历失败时，函数 MUST 降级为安全的字节截断策略，仍返回不超过预算的结果并设置截断标志，不得 panic。`TruncateMarkdown` 降级截断 MUST NOT 追加后缀；后缀追加仅由 `TruncateMarkdownWithSuffix` 负责。

#### Scenario: 解析失败降级

- **WHEN** goldmark 解析或 AST 遍历返回错误
- **THEN** `TruncateMarkdown` 按 UTF-8 安全边界做简单字节截断，不追加后缀，截断标志为 `true`
