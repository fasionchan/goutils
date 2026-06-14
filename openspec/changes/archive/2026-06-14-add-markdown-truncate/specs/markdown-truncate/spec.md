## ADDED Requirements

### Requirement: 按字节上限截断 Markdown 内容

`TruncateMarkdown` SHALL 接受 Markdown 字符串与 `maxBytes` 上限，返回截断后的字符串及是否发生截断的布尔值。字节计数 MUST 按 UTF-8 编码计算，与 Go `len([]byte(s))` 一致。

#### Scenario: 内容未超出上限

- **WHEN** 输入 Markdown 的 UTF-8 字节长度小于或等于 `maxBytes`
- **THEN** 返回原始内容不变，截断标志为 `false`

#### Scenario: 内容超出上限

- **WHEN** 输入 Markdown 的 UTF-8 字节长度大于 `maxBytes`
- **THEN** 返回长度不超过 `maxBytes` 的 Markdown 字符串，截断标志为 `true`

#### Scenario: maxBytes 非正数

- **WHEN** `maxBytes` 小于或等于 0
- **THEN** 返回空字符串，截断标志为 `true`（若原内容非空）或 `false`（若原内容为空）

### Requirement: 在块级节点边界截断

截断 MUST 在 goldmark AST 块级节点边界进行，不得将单个块级节点从中间切开。若某块级节点本身超过 `maxBytes`（含省略标记预留空间），该节点 MUST 作为整体保留或整体丢弃，不得产生不完整的 Markdown 结构。

#### Scenario: 多段落截断

- **WHEN** Markdown 包含多个段落且总长度超出 `maxBytes`
- **THEN** 结果包含从文档开头起、完整保留的连续块级节点，不包含被截断节点的不完整片段

#### Scenario: 代码块不截断中间

- **WHEN** Markdown 包含 fenced code block 且截断点落在该代码块内
- **THEN** 要么完整保留该代码块（若剩余空间足够），要么不包含该代码块，不得输出未闭合的 fence 标记

#### Scenario: 列表项保持完整

- **WHEN** Markdown 包含有序或无序列表且截断点落在某列表项内
- **THEN** 该列表项 MUST 作为整体保留或整体丢弃

### Requirement: 截断时追加省略标记

当发生截断时，函数 SHALL 在结果末尾追加省略标记以提示内容被截断。默认省略标记为 `...`，其字节长度 MUST 计入 `maxBytes` 预算。

#### Scenario: 默认省略标记

- **WHEN** 内容被截断且未指定自定义省略标记
- **THEN** 结果以 `...` 结尾，且总字节长度不超过 `maxBytes`

#### Scenario: 省略标记占用预算

- **WHEN** `maxBytes` 较小，扣除省略标记后无法容纳任何块级节点
- **THEN** 返回仅含省略标记（或空字符串，若 `maxBytes` 不足以容纳省略标记）的结果，截断标志为 `true`

### Requirement: 空输入与纯空白输入

#### Scenario: 空字符串

- **WHEN** 输入为空字符串
- **THEN** 返回空字符串，截断标志为 `false`

#### Scenario: 仅空白字符

- **WHEN** 输入仅包含空白字符且未超出 `maxBytes`
- **THEN** 原样返回，截断标志为 `false`

### Requirement: AST 解析失败时的降级策略

当 goldmark AST 解析或遍历失败时，函数 MUST 降级为安全的字节截断策略，仍返回不超过 `maxBytes` 的结果并设置截断标志，不得 panic。

#### Scenario: 解析失败降级

- **WHEN** goldmark 解析或 AST 遍历返回错误
- **THEN** 按 UTF-8 安全边界做简单字节截断（不截断多字节字符中间），追加省略标记，截断标志为 `true`
