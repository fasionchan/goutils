## Why

goutils 的 `mdutils` 包已提供 `SplitMarkdown`，可将 Markdown 按字节上限拆成多段，但缺少「截断为单段」的能力。在企微 Markdown 消息（2048 字节上限）、日志预览、通知摘要等场景，需要的是保留语法结构、在边界内尽量多保留内容的单段截断，而非拆分。补齐该能力可复用现有 goldmark AST 解析思路，与 `SplitMarkdown` 形成配套工具集。

## What Changes

- 在 `baseutils/mdutils` 新增 `TruncateMarkdown(content string, maxBytes int) (string, bool)` 函数
- 基于 goldmark AST 在块级节点边界截断，避免截断代码块、列表、标题等结构中间
- 超出上限时追加可配置的省略标记（默认 `...`），并返回 `truncated=true`
- 内容未超出上限时原样返回，`truncated=false`
- 补充单元测试，覆盖常见 Markdown 结构与边界情况
- 不修改现有 `SplitMarkdown` 的 API 与行为

## Capabilities

### New Capabilities

- `markdown-truncate`: Markdown 内容按字节上限截断，保持块级语法结构完整，并指示是否发生截断

### Modified Capabilities

（无）

## Impact

- **代码**: `baseutils/mdutils/` 新增 `truncate.go`、`truncate_test.go`
- **依赖**: 复用现有 `github.com/yuin/goldmark`，无新增外部依赖
- **API**: 纯新增函数，无 breaking change
- **调用方**: 可被 opsys/server 等项目的企微 Markdown 消息、日志/通知摘要等场景引用
