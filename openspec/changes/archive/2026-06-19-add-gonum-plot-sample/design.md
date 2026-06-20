## Context

gonum/plot 渲染折线图需要 `plotter.XYer`（如 `plotter.XYs`）及 `plot.Plot` 配置（Title、Legend、Axis 等）。监控场景常见数据来源包括 Prometheus，但 **Prometheus 响应解析与 PromQL 不属于 gonumplot 职责**，将在独立包中实现并产出 `Chart`。

gonumplot 定位为：**Chart 数据结构 + 渲染 + 可运行的监控类示例**，供上层（含未来 Prometheus 包）复用。

## Goals / Non-Goals

**Goals:**

- 定义可 JSON 序列化的 `Chart` / `Series` / `Point`
- 实现 `RenderChart` 将 `Chart` 渲染为 PNG
- 提供 CPU、内存、文件系统、磁盘 I/O、网络流量 5 类示例（Chart JSON fixture + 渲染函数）
- 单元测试覆盖 JSON 往返与渲染

**Non-Goals:**

- **不实现** Prometheus HTTP 客户端、PromQL 查询、Prometheus JSON → Chart 转换（其他包负责）
- 不做交互式 Web 图表（仅静态 PNG）
- 不集成 Grafana

## Decisions

### 1. Chart 作为唯一输入模型

**选择**: 渲染层只接受 `Chart`；示例从 Chart JSON fixture 加载（`json.Unmarshal`）

```go
type Chart struct {
    Title   string   `json:"title"`
    XLabel  string   `json:"xLabel"`
    YLabel  string   `json:"yLabel"`
    Series  []Series `json:"series"`
}

type Series struct {
    Name   string  `json:"name"`
    Points []Point `json:"points"`
}

type Point struct {
    X float64 `json:"x"` // Unix 时间戳（秒）
    Y float64 `json:"y"`
}
```

**理由**:
- JSON 可序列化，fixture  human-readable，测试简单
- 与 Prometheus 解耦；未来 `libs/promplot`（或类似包）负责 `PrometheusResponse → Chart`
- 渲染代码稳定，不随数据源变化

**与先前方案的变更**: 不再在 gonumplot 内做 Prometheus 转换；fixture 存 **Chart JSON**，非 PromQL API 响应。

### 2. 包结构

```
libs/gonumplot/
  model.go              // Chart, Series, Point；LoadChartJSON
  render.go             // RenderChart, RenderChartPNG
  examples/
    cpu.go
    memory.go
    filesystem.go       // 文件系统空间/挂载点等
    diskio.go           // 磁盘读写字节率
    network.go
  testdata/
    cpu.json            // Chart JSON
    memory.json
    filesystem.json
    diskio.json
    network.json
  *_test.go
```

### 3. 示例数据语义（Chart fixture 内容）

| 示例 | 典型监控含义 | Y 轴单位 | Series 示例 |
|------|-------------|---------|-------------|
| CPU | 使用率或多核 load | `%` 或 cores | `cpu0`, `cpu1` |
| Memory | 已用/可用内存 | `bytes` 或 `%` | `used`, `cached` |
| Filesystem | 挂载点空间使用 | `bytes` 或 `%` | `/`, `/data` |
| Disk I/O | 读/写吞吐 | `bytes/s` | `read`, `write` |
| Network | 收/发流量 | `bytes/s` | `rx`, `tx` |

fixture 为合理模拟时序点（Unix X + float Y），无需真实 Prometheus 数值。

### 4. 图表类型

**选择**: 5 类示例均用折线图（`plotter.NewLine`）；多 series 共用 Legend

### 5. 时间轴

**选择**: `Point.X` 为 Unix 秒；render 时 X 轴按 gonum/plot 默认 tick 格式化（首版不自定义时区）

## Risks / Trade-offs

- **[Trade-off] 示例 fixture 需手写 Chart JSON** → 比 PromQL 响应冗长，但职责清晰；Prometheus 包可自动生成 Chart
- **[Risk] gonum/plot 中文标题** → Title/YLabel UTF-8 字符串，依赖默认字体
- **[Risk] vg/png 依赖** → gonum 标准做法

## Migration Plan

纯新增包，无迁移。

## Open Questions

- Prometheus 转换包命名与路径（如 `libs/promplot`）留待后续 change，本变更不创建
