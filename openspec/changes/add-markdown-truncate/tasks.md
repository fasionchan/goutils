## 1. 基础实现

- [x] 1.1 在 `baseutils/mdutils/truncate.go` 实现 `TruncateMarkdown(content string, maxBytes int) (string, bool)`，默认 suffix 为 `...`
- [x] 1.2 实现 `TruncateMarkdownWithSuffix(content string, maxBytes int, suffix string) (string, bool)`，suffix 字节计入 maxBytes 预算
- [x] 1.3 复用 goldmark AST 遍历逻辑，按块级节点顺序累加源码，在超出预算前停止
- [x] 1.4 实现 `maxBytes <= 0`、空字符串、纯空白输入的边界处理
- [x] 1.5 实现 AST 解析失败时的 UTF-8 安全字节截断降级函数 `fallbackTruncate`

## 2. 超大节点与结构完整性

- [x] 2.1 处理单个块级节点超过剩余预算的情况：有已累加内容则丢弃该节点；无已累加内容则 UTF-8 安全截断该节点
- [x] 2.2 确保代码块（fenced code block）不被从 fence 中间截断
- [x] 2.3 确保列表项作为整体保留或丢弃

## 3. 单元测试

- [x] 3.1 添加 `truncate_test.go`：内容未超出上限时原样返回
- [x] 3.2 测试多段落截断，验证只保留完整块级节点
- [x] 3.3 测试含 fenced code block 的 Markdown 截断
- [x] 3.4 测试含有序/无序列表的 Markdown 截断
- [x] 3.5 测试超大单节点降级截断
- [x] 3.6 测试 `maxBytes <= 0`、空输入、自定义 suffix
- [x] 3.7 运行 `go test ./baseutils/mdutils/...` 确保全部通过

## 4. 收尾

- [x] 4.1 确认未修改 `SplitMarkdown` 现有行为与测试
- [x] 4.2 为公开函数补充简要 Go doc 注释（中文说明用途与返回值含义）
