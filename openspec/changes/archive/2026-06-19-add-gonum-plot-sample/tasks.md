## 1. 数据模型与依赖

- [x] 1.1 在 `libs/gonumplot/model.go` 定义 `Chart`、`Series`、`Point` 及 JSON tag
- [x] 1.2 实现 `LoadChartJSON(r io.Reader) (*Chart, error)` 与 `MarshalJSON` 友好结构
- [x] 1.3 添加 `gonum.org/v1/plot` 依赖（`go get gonum.org/v1/plot`）
- [x] 1.4 编写 `model_test.go`：Chart JSON 往返、LoadChartJSON 测试

## 2. 渲染层

- [x] 2.1 在 `render.go` 实现 `RenderChart(chart *Chart, w io.Writer) error` 与 `RenderChartPNG(chart *Chart) ([]byte, error)`
- [x] 2.2 使用 `plot.New`、`plotter.NewLine`、Legend、Title/X/Y 标签
- [x] 2.3 空 series 返回明确错误
- [x] 2.4 编写 `render_test.go`：PNG smoke test（魔数、非空）

## 3. Chart JSON Fixture

- [x] 3.1 在 `libs/gonumplot/testdata/` 添加 5 个 Chart JSON：`cpu.json`、`memory.json`、`filesystem.json`、`diskio.json`、`network.json`
- [x] 3.2 每个 fixture 含合理时序点、标题、轴标签及多 series（如需要）

## 4. 监控示例

- [x] 4.1 `examples/cpu.go`：加载 cpu.json → 渲染 PNG
- [x] 4.2 `examples/memory.go`：内存示例
- [x] 4.3 `examples/filesystem.go`：文件系统空间示例
- [x] 4.4 `examples/diskio.go`：磁盘读/写 I/O 示例
- [x] 4.5 `examples/network.go`：网络 rx/tx 示例
- [x] 4.6 每个 example 编写 smoke test

## 5. 收尾

- [x] 5.1 为公开类型与函数补充中文 Go doc 注释
- [x] 5.2 运行 `go test ./libs/gonumplot/...` 确保全部通过
- [x] 5.3 确认包内无 Prometheus 相关代码
