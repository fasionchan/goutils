## ADDED Requirements

### Requirement: Chart 数据模型与 JSON 序列化

包 SHALL 定义 `Chart`、`Series`、`Point` 结构，并通过 `json` tag 支持标准 JSON 序列化与反序列化。`Chart` MUST 包含 `title`、`xLabel`、`yLabel`、`series`；每个 `Series` MUST 包含 `name` 与 `points`（X/Y 为 float64）。

#### Scenario: Chart JSON 往返

- **WHEN** 构造含多条 series 的 `Chart` 并 `json.Marshal` / `json.Unmarshal`
- **THEN** 字段值与 series 数量、点坐标保持一致

#### Scenario: 从 JSON 文件加载

- **WHEN** 调用 `LoadChartJSON(r io.Reader)` 或等价函数读取合法 Chart JSON
- **THEN** 返回 `*Chart`，无 error

#### Scenario: 空 Chart

- **WHEN** `Chart` 不含任何 series
- **THEN** 渲染函数返回明确错误，不 panic

### Requirement: Chart 渲染为 PNG

包 SHALL 提供将 `Chart` 渲染为 PNG 的函数，使用 gonum/plot 折线图（`plotter.Line`）。MUST 支持写入 `io.Writer` 或返回 `[]byte`。

#### Scenario: 单 series 折线图

- **WHEN** `Chart` 含一条 series 且至少 2 个点
- **THEN** 输出非空 PNG，魔数为 `\x89PNG`

#### Scenario: 多 series 与图例

- **WHEN** `Chart` 含多条 series
- **THEN** 成功渲染，含 Legend

#### Scenario: 轴标签与标题

- **WHEN** `Chart` 设置了 Title、XLabel、YLabel
- **THEN** 渲染不报错，配置写入 plot

### Requirement: CPU 监控示例

包 SHALL 提供 CPU 使用率（或多核）时序折线图示例，从 Chart JSON fixture 加载并渲染 PNG。

#### Scenario: CPU 示例端到端

- **WHEN** 调用 CPU 示例渲染函数
- **THEN** 成功输出 PNG；fixture 含 CPU 相关 series 名称与 `%` 或等价 Y 轴标签

### Requirement: 内存监控示例

包 SHALL 提供内存使用量时序折线图示例，从 Chart JSON fixture 加载并渲染。

#### Scenario: 内存示例端到端

- **WHEN** 调用内存示例渲染函数
- **THEN** 成功输出 PNG；Y 轴标签体现 bytes 或百分比

### Requirement: 文件系统监控示例

包 SHALL 提供文件系统空间/挂载点时序折线图示例（如各挂载点使用率），从 Chart JSON fixture 加载并渲染。

#### Scenario: 文件系统示例端到端

- **WHEN** 调用文件系统示例渲染函数
- **THEN** 成功输出 PNG；可含多个挂载点 series

### Requirement: 磁盘 I/O 监控示例

包 SHALL 提供磁盘读/写吞吐时序折线图示例，从 Chart JSON fixture 加载并渲染。

#### Scenario: 磁盘 I/O 示例端到端

- **WHEN** 调用磁盘 I/O 示例渲染函数
- **THEN** 成功输出 PNG；read/write 以不同 series 区分

### Requirement: 网络流量监控示例

包 SHALL 提供网络接收/发送流量时序折线图示例，从 Chart JSON fixture 加载并渲染。

#### Scenario: 网络流量示例端到端

- **WHEN** 调用网络示例渲染函数
- **THEN** 成功输出 PNG；rx/tx 以不同 series 区分

### Requirement: 测试不依赖外部服务

单元测试 MUST 使用内嵌或 `testdata/` Chart JSON fixture，不依赖 Prometheus 或其他外部服务。

#### Scenario: 模型与渲染测试

- **WHEN** 运行 `go test ./libs/gonumplot/...`
- **THEN** JSON 往返与 PNG smoke test 全部通过

### Requirement: Prometheus 逻辑不在本包

gonumplot 包 MUST NOT 包含 Prometheus Query API 解析、PromQL 客户端或 Prometheus JSON 到 Chart 的转换。上述能力由其他包后续提供，其输出 MUST 为与本包兼容的 `Chart` 结构。

#### Scenario: 包边界

- **WHEN** 检查 `libs/gonumplot/` 源码
- **THEN** 不存在 `prometheus`、`promql`、`matrix`/`vector` resultType 等 Prometheus 专用转换代码
