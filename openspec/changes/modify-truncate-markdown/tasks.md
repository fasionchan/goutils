## 1. 重构 API

- [ ] 1.1 将截断核心逻辑移入 `TruncateMarkdown`，`maxBytes` 全部用于内容，不追加后缀
- [ ] 1.2 重写 `TruncateMarkdownWithSuffix`：调用 `TruncateMarkdown`，在截断结果上追加后缀
- [ ] 1.3 `TruncateMarkdownWithSuffix` 中 `suffix` 为空时使用默认后缀 `...`
- [ ] 1.4 更新两个公开函数的 Go doc 注释，明确职责划分

## 2. 更新测试

- [ ] 2.1 更新 `TruncateMarkdown` 相关测试：截断结果不包含 `...` 后缀
- [ ] 2.2 新增或调整 `TruncateMarkdownWithSuffix` 测试：空 suffix 使用默认后缀、自定义 suffix
- [ ] 2.3 确认 `TruncateMarkdownWithSuffix` 总输出不超过 `maxBytes`
- [ ] 2.4 运行 `go test ./baseutils/mdutils/...` 确保全部通过

## 3. 收尾

- [ ] 3.1 确认 `TruncateMarkdownWithSuffix` 不再包含独立截断逻辑，仅委托 `TruncateMarkdown`
