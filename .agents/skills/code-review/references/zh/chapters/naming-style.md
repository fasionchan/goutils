# 命名与风格

通用 Go 命名与代码风格检查项（参考 Effective Go / Go Code Review Comments 精神）。

## 检查项

1. **使用 Go 惯用命名**：导出标识用 `MixedCaps`，包内可用 `mixedCaps`；避免 `snake_case` 作标识符名（测试文件名除外）。
2. **避免无意义缩写与冗余类型后缀**：如 `userMap map[string]*User` 可接受，但 `GetUserInfo` 若返回已是 User 则勿堆砌 `Info/Data/Object`；缩写保持全大写或全小写一致（`URL`/`url`，勿 `Url`）。
3. **包名简短、小写、无下划线**：包名应为单数名词或简短词，避免与标准库常见包名无必要冲突，且勿在包内重复包名（如 `http.HTTPClient`）。
4. **格式与风格工具一致**：变更应可通过 `gofmt` / `goimports` 风格检查；勿手工保留明显错位缩进或杂乱 import 分组（除非生成代码且有明确标记）。
5. **错误变量与布尔命名清晰**：错误值优先 `err`；布尔名体现语义（`isReady`/`hasCache`），避免双重否定命名。
