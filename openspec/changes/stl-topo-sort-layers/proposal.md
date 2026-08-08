## Why

`stl` 包已提供一维拓扑排序（`Graph.TopoSort` / `TopoSortDataByFormers`），可得到合法的全局顺序，但无法直接表达「同一批互不依赖、可并行处理」的节点集合。需要在保留现有 API 的前提下，增加基于 Kahn 算法的分层拓扑排序，供调度、批处理等场景按层消费。

## What Changes

- 在 `stl.Graph[T]` 上新增 `TopoSortLayers`，返回 `[][]T`（每层为一批当前入度为 0 的节点）
- 新增 `TopoSortDataByFormersLayers`，与现有 `TopoSortDataByFormers` 参数对称，返回 `[][]Data`
- 层内节点顺序沿用传入的 `container` 语义；`container == nil` 时默认栈，与现有一维 API 一致
- 有环、仅作为边目标出现的节点、未知 former 等边界行为与现有一维拓扑排序**完全对齐**
- **不修改**现有 `TopoSort` / `TopoSortDataByFormers` 的签名与语义
- 补充覆盖分层语义、边界与 container 顺序的单元测试

## Capabilities

### New Capabilities

- `stl-topo-sort-layers`: `stl` 包 Kahn 分层拓扑排序 API（Graph 与 DataByFormers 两侧），定义层语义、层内顺序、环/未知 former 对齐行为与可观测验收条件

### Modified Capabilities

（无）

## Impact

- **代码**: `stl/graph.go` 新增两个公开函数；`stl/graph_test.go` 新增分层相关测试
- **依赖**: 无新增外部依赖；复用现有 `Container` / heap / stack
- **API**: 纯新增，无 breaking change
- **调用方**: 需要「按可并行批」处理依赖图的代码可选用新 API；现有一维调用方不受影响
