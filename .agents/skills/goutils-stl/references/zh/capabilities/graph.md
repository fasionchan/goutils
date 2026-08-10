# 能力：graph

源码重点：`stl/graph.go`，容器见 `stl/container.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 邻接表 | `Graph`、`GraphFromFormers` |
| 拓扑序 | `Graph.TopoSort`、`Graph.TopoSortLayers` |
| 数据级排序 | `TopoSortDataByFormers`、`TopoSortDataByFormersLayers` |
| 就绪集容器 | `NewStackAsContainer`、`NewMinHeapAsContainer`、`NewMaxHeapAsContainer` |

容器工厂传 `nil` 使用默认栈序；需要按 key 决胜时再传入堆工厂。
