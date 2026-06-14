## Context

`baseutils/mdutils` 已实现 `TruncateMarkdown` 与 `TruncateMarkdownWithSuffix`，但当前 `TruncateMarkdown` 委托后者并默认追加 `...`，导致基础 API 无法提供纯截断。用户要求反转调用关系：`TruncateMarkdown` 负责纯截断，`TruncateMarkdownWithSuffix` 在其基础上追加后缀。

## Goals / Non-Goals

**Goals:**

- `TruncateMarkdown(content, maxBytes)` 仅截断，`maxBytes` 全部用于内容
- `TruncateMarkdownWithSuffix(content, maxBytes, suffix)` 调用 `TruncateMarkdown`，追加后缀；`suffix == ""` 时使用 `...`
- 将现有截断核心逻辑（边界处理、AST 遍历、降级）收敛到 `TruncateMarkdown`
- 更新测试与文档

**Non-Goals:**

- 不改变 AST 截断算法本身
- 不新增其他 API 变体
- 不修改 `SplitMarkdown`

## Decisions

### 1. 调用关系反转

**选择**: `TruncateMarkdownWithSuffix` → `TruncateMarkdown`

**理由**: 符合用户要求，职责单一——截断与后缀追加分离。

**备选**: 共享内部函数 `truncateMarkdownCore` —— 过度抽象，两层公开 API 已足够。

### 2. TruncateMarkdown 实现

**选择**: 将当前 `TruncateMarkdownWithSuffix` 中的截断逻辑（不含 suffix 处理）移入 `TruncateMarkdown`，`contentBudget = maxBytes`。

```go
func TruncateMarkdown(content string, maxBytes int) (string, bool) {
    // 边界处理 + truncateMarkdownByAST / fallbackTruncate
    // 不追加后缀
}

func TruncateMarkdownWithSuffix(content string, maxBytes int, suffix string) (string, bool) {
    if suffix == "" {
        suffix = defaultTruncateSuffix
    }
    if len(content) <= maxBytes {
        return content, false
    }
    contentBudget := maxBytes - len(suffix)
    if contentBudget < 0 {
        contentBudget = 0
    }
    truncated, ok := TruncateMarkdown(content, contentBudget)
    // 追加 suffix 逻辑
}
```

### 3. suffix 为空时使用默认后缀

**选择**: `TruncateMarkdownWithSuffix` 内部 `if suffix == "" { suffix = "..." }`

**理由**: 用户明确要求；空 suffix 不再表示「无后缀」，而是「使用默认后缀」。需要无后缀的调用方应直接使用 `TruncateMarkdown`。

### 4. 迁移策略

**BREAKING**: 原 `TruncateMarkdown` 调用方若依赖默认 `...` 后缀，需改为 `TruncateMarkdownWithSuffix(content, maxBytes, "")`。

## Risks / Trade-offs

- **[Risk] BREAKING 变更影响已有调用方** → 当前 goutils 内部无其他引用；发布时在 changelog 标注
- **[Risk] suffix 语义变化（空=默认而非无后缀）** → 文档明确说明；纯截断用 `TruncateMarkdown`

## Migration Plan

1. 重构 `truncate.go` 调用关系
2. 更新 `truncate_test.go`：`TruncateMarkdown` 测试不期望后缀；后缀相关测试改用 `TruncateMarkdownWithSuffix`
3. 运行 `go test ./baseutils/mdutils/...`

## Open Questions

（无）
