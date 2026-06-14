## Why

当前 `TruncateMarkdown` 内部委托 `TruncateMarkdownWithSuffix` 并默认追加 `...` 后缀，职责边界不清晰：调用方若只需纯截断（如预留空间给外部拼接），无法通过基础 API 实现。应将「截断内容」与「追加省略标记」拆分为两层，使 API 语义更明确、组合更灵活。

## What Changes

- **BREAKING**: `TruncateMarkdown` 仅执行截断，不追加任何后缀；`maxBytes` 全部用于内容预算
- `TruncateMarkdownWithSuffix` 改为调用 `TruncateMarkdown`，在截断结果上追加后缀
- `TruncateMarkdownWithSuffix` 的 `suffix` 为空字符串时，使用默认后缀 `...`
- 后缀字节长度仍计入 `maxBytes` 总预算（由 `TruncateMarkdownWithSuffix` 从预算中扣除后再调用 `TruncateMarkdown`）
- 更新单元测试以反映新 API 语义
- 更新 Go doc 注释

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `markdown-truncate`: 调整 `TruncateMarkdown` 与 `TruncateMarkdownWithSuffix` 的职责划分与调用关系

## Impact

- **代码**: `baseutils/mdutils/truncate.go`、`truncate_test.go`
- **API**: `TruncateMarkdown` 行为变更（**BREAKING**），不再默认追加 `...`
- **调用方**: 需要省略标记的场景应改用 `TruncateMarkdownWithSuffix(content, maxBytes, "")` 或显式传入后缀
- **依赖**: 无新增外部依赖
