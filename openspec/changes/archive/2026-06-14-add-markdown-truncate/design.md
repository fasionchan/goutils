## Context

`baseutils/mdutils` 已有 `SplitMarkdown`，通过 goldmark 解析 AST、按块级节点拆分多段 Markdown。与之相对，截断场景（企微 Markdown 2048 字节限制、日志/通知摘要）需要单段输出，并在语法结构边界处截断。

现有 `SplitMarkdown` 的实现（`.split.go`）已验证 goldmark AST 遍历、块级节点 `Lines().Pos()` 提取源码片段的模式，可复用于截断逻辑。

## Goals / Non-Goals

**Goals:**

- 提供 `TruncateMarkdown(content string, maxBytes int) (string, bool)` 及可选的 `TruncateMarkdownWithSuffix(content string, maxBytes int, suffix string) (string, bool)`
- 在块级 AST 节点边界累加内容，超出 `maxBytes`（含 suffix 预算）时停止
- 与 `SplitMarkdown` 保持一致的 goldmark 依赖与代码风格
- 覆盖单元测试：短文本、多段落、代码块、列表、超大单节点、解析降级

**Non-Goals:**

- 不修改 `SplitMarkdown` 行为
- 不实现 inline 级别精细截断（如句子、词语边界）
- 不处理 HTML 渲染或 Markdown 转纯文本
- 不在本变更中集成到 opsys/server 调用方

## Decisions

### 1. API 设计：返回值 `(string, bool)` 而非结构体

**选择**: 简单元组 `(result, truncated bool)`

**理由**: 与 Go 标准库 `strings.Cut` 等惯例一致，调用方多数场景只需判断是否截断。若后续需要 `totalBytes` 等元数据，可再扩展 `TruncateResult` 结构体。

**备选**: 返回 `TruncateResult{Content, Truncated, OriginalBytes}` —— 过度设计，当前需求不需要。

### 2. 截断算法：AST 块级节点累加

**选择**: 复用 `SplitMarkdown` 的 AST 遍历模式，按文档顺序累加块级节点源码，直到再加一个节点会超出 `maxBytes - len(suffix)` 为止。

**理由**: 与现有 `SplitMarkdown` 一致，保证块级结构完整；实现成本低。

**备选**:
- 纯字节截断 —— 会破坏 Markdown 语法，不可接受
- 先 `SplitMarkdown` 再取首段 —— 语义不同（首段可能仍超限），且浪费计算

### 3. 超大单节点处理：整体保留或整体丢弃

**选择**: 若单个块级节点字节数 > `maxBytes - len(suffix)`，则丢弃该节点及之后所有内容，返回已累加部分 + suffix；若之前无任何内容，则尝试 UTF-8 安全字节截断该节点本身。

**理由**: 与 spec 要求一致；超大代码块场景下至少保留前缀内容比返回空更有用。

### 4. 省略标记：独立函数 + 默认 `...`

**选择**: `TruncateMarkdown` 调用 `TruncateMarkdownWithSuffix(content, maxBytes, "...")`

**理由**: 便于调用方自定义（如 `…` 或 `\n...(已截断)`），默认行为满足大多数场景。

### 5. 降级策略：UTF-8 安全字节截断

**选择**: AST 失败时，使用 `[]byte` 截断并回退至 rune 边界，再追加 suffix。

**理由**: 与 spec 一致，保证函数永不 panic、始终有界输出。

## Risks / Trade-offs

- **[Risk] 块级截断可能丢失大量尾部内容** → 接受；这是保持语法完整的必要代价，调用方可增大 `maxBytes` 或使用 `SplitMarkdown` 分段发送
- **[Risk] 省略标记可能与 Markdown 语法冲突** → suffix 默认 `...` 作为纯文本追加，不解析为语法；文档说明自定义 suffix 需自行确保合法
- **[Risk] goldmark 版本升级可能改变 AST 边界** → 单元测试锁定行为；与 `SplitMarkdown` 共享同一依赖版本
- **[Trade-off] 不做 inline 截断** → 短 `maxBytes` 下可能只保留很少内容；符合 Non-Goals

## Migration Plan

纯新增 API，无需迁移。发布新版本 goutils 后，调用方可按需引入 `mdutils.TruncateMarkdown`。

## Open Questions

- 是否在首版暴露 `TruncateMarkdownWithSuffix`，还是仅内部使用？建议首版公开，成本低且灵活。
- 超大单节点降级截断时，是否优先在换行符处截断？可作为后续优化，首版仅做 UTF-8 安全截断。
