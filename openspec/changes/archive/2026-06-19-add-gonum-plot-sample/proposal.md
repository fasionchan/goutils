## Why

goutils 目前缺少基于 Go 原生绘图库的可复用监控图表示例。gonum/plot 的数据结构（`plotter.XYs` 等）与业务侧时序数据并不直接对应，需要一层可 JSON 序列化的 `Chart` 模型及统一渲染逻辑。Prometheus PromQL 转换将在**其他包**后续实现；本次 `libs/gonumplot` 仅负责 Chart 定义、画图与监控类示例。

## What Changes

- 在 `libs/gonumplot/` 新增包，基于 `gonum.org/v1/plot` 实现折线图渲染
- 定义 `Chart` / `Series` / `Point` 数据模型，**支持 JSON 序列化/反序列化**
- 实现 `RenderChart` 将 `Chart` 渲染为 PNG
- 提供 5 类监控指标示例：CPU、内存、文件系统（空间/挂载）、磁盘 I/O、网络流量
- 示例使用 **Chart JSON fixture**（非 Prometheus 响应格式），演示从结构化数据到 PNG 的完整路径
- 补充单元测试：模型 JSON 往返、渲染 smoke test（不依赖外部服务）

## Capabilities

### New Capabilities

- `gonumplot-chart`: 可 JSON 序列化的 Chart 数据模型、gonum/plot 渲染及监控指标示例

### Modified Capabilities

（无）

## Impact

- **代码**: 新增 `libs/gonumplot/`（model、render、examples、testdata）
- **依赖**: 新增 `gonum.org/v1/plot`（及 `vg/png` 等子包）
- **API**: 纯新增；`Chart` JSON schema 供后续 Prometheus 转换包及其他数据源复用
- **不在本包**: Prometheus Query API 解析、PromQL 客户端、metric labels 合成
