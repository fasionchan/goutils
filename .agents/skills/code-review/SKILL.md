---
name: code-review
description: >-
  Go 代码评审框架（风格/惯例/轻量安全）。在用户要求 code review、评审 PR/diff、
  或检查 Go 代码风格与惯例时加载。非 Go 为主或与评审无关的任务不要使用。
---

# Go code-review skill

通用 Golang 风格 / 规范 / 惯例评审框架。首版只含通用检查项；项目个性化原则放在扩展位，扩展位为空时不得编造私有规则。

## When to use

在以下场景**按需**加载本 Skill（改完代码准备提交时**不**强制加载）：

1. 用户要求 **code review** / 代码评审
2. 用户要求评审 **PR** 或 **diff** / 变更
3. 用户要求检查 **Go 代码风格与惯例**

评审对象默认为当前变更或用户指定路径中的 Go 代码；首版不对接 GitHub PR API。

## When NOT to use / Out of scope

- 变更以**非 Go** 源码为主（标注「本 Skill 不适用」并跳过或降级说明，不强行套用）
- 与代码评审无关的任务（实现新功能、写文档、跑 CI 等）
- **完整安全审计**、**性能 profiling**、**业务正确性证明** — 不在首版强制范围（见 `references/zh/OVERVIEW.md` 与 `security-basics` 章）

## Language gate（mandatory）

首版**只加载** `references/zh/`。

| 偏好 | 加载 |
|------|------|
| 中文 / 默认（首版） | `references/zh/` |
| 英文 | `references/en/`（首版仅为占位，**不要**加载其实质正文） |

同一会话不要混开另一语言树的正文。

## Progressive loading

1. 读本文件（触发、边界、输出模板）
2. 打开 `references/zh/INDEX.md`（章目录与加载顺序）
3. 浏览 `references/zh/OVERVIEW.md`（范围、严重度、非范围）
4. 按 INDEX 顺序打开 `references/zh/chapters/*.md` 通用章
5. 若 `references/zh/project-specific/` 有**实质**个性化内容则再加载；若仅占位则**跳过**且**不编造**项目私有规则

## 扩展加载顺序（AC-9）

1. **先**通用章（`chapters/`）
2. **再**扩展位（`project-specific/`）
3. 扩展位为空或仅占位说明 → 跳过，不得编造 goutils 或其它项目私有规则

## Review output template

评审回复**必须**使用以下结构（标题可微调，字段语义不可缺）。无某级问题时保留标题并写「无」。

```markdown
## 评审结论摘要
- 总体结论：通过 | 有条件通过 | 不通过
- 范围说明：审查了哪些路径/文件；跳过了哪些（含「非 Go」）
- 统计：blocker=N, major=N, minor=N, nit=N

## 问题列表（按严重度）

### blocker
- **[标题]**
  - 位置：`path/to/file.go:行`（或 diff hunk 标识）
  - 证据：简述违反的检查项/现象
  - 建议：可执行的修改方向

### major
（同上条目结构）

### minor
（同上）

### nit
（同上）

## 未覆盖 / 超出范围
- （完整安全审计、性能 profiling、业务正确性等 — 显式声明未做）
```

严重度定义见 `references/zh/OVERVIEW.md`：`blocker` / `major` / `minor` / `nit`。

## Examples

人工对照示例（故意含缺陷）：

- `examples/flawed-snippet.go` — 含命名与错误处理等问题
- `examples/expected-findings.md` — 预期命中要点（≥2 类）

实现或自检时可对示例走查一遍，确认检查项与输出模板可映射。
