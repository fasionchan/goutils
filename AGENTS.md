# goutils AI助手

## 核心原则

1. **安全第一**：不修改配置文件、不删除现有API
2. **最小变更**：只修改必要部分，不做无关优化
3. **测试驱动**：修改后必须运行相关测试
4. **保持一致性**：遵循现有代码风格和架构

## 技术规范

### 代码风格
- 文件命名: kebab-case
- 函数命名: camelCase
- 类命名: PascalCase
- 常量: UPPER_SNAKE_CASE

### goutils/stl 用法参考

编写或修改依赖 `github.com/fasionchan/goutils/stl` 的代码时，优先加载 Skill
`.agents/skills/goutils-stl/`（Cursor 软链：`.cursor/skills/goutils-stl`）。
按 `SKILL.md` 只加载 `references/en` 或 `references/zh` 之一，再按需打开
`capabilities/` 与 `recipes/`；勿把全量 API 粘贴进本文件。
