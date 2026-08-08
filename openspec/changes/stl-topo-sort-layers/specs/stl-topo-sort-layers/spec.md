## ADDED Requirements

### Requirement: Graph 侧 Kahn 分层拓扑排序

`Graph[T].TopoSortLayers` SHALL 接受与 `TopoSort` 相同的 `container` 工厂参数，返回 `[][]T`。每一层 MUST 是 Kahn 算法中同一波次内入度为 0 的节点集合：处理完本层全部节点的出边后，新产生的零入度节点属于下一层。

#### Scenario: 无环图满足层间拓扑与层内无依赖

- **WHEN** 对无环有向图调用 `TopoSortLayers`
- **THEN** 返回的层序列中，对图中任意边 `u→v`（且 `u`、`v` 均出现在结果中），`u` 所在层号严格小于 `v` 所在层号；同一层内任意两节点之间不存在依赖边

#### Scenario: 空图

- **WHEN** 图为空（无节点）
- **THEN** 返回空的分层结果（`assert.Empty` 意义下为空；`nil` 或长度 0 的切片均可）

#### Scenario: 单节点

- **WHEN** 图仅含一个节点且无出边
- **THEN** 返回恰好一层，该层仅含该节点

#### Scenario: 多源汇入单目标

- **WHEN** 图为 `a→c`、`b→c`（例如 `{"a":{"c"},"b":{"c"},"c":nil}`）且使用 `NewMinHeapAsContainer`
- **THEN** 第一层为 `["a","b"]`（堆序），第二层为 `["c"]`

### Requirement: 层内顺序沿用 container

层内节点排列顺序 MUST 由传入的 `container` 的连续 `Pop` 顺序决定。当 `container == nil` 时，MUST 默认使用 `NewStackAsContainer`，与 `TopoSort` 一致。

#### Scenario: MinHeap 层内有序

- **WHEN** 同一层存在多个初始零入度节点且传入 `NewMinHeapAsContainer`
- **THEN** 该层内节点按升序排列（与同参数下一维 `TopoSort` 对这些节点的相对出队顺序一致）

#### Scenario: 默认栈可用且结果合法

- **WHEN** `container` 为 `nil`
- **THEN** 分层结果仍满足层间拓扑与层内无依赖；不要求与 map 迭代相关的绝对顺序稳定

### Requirement: 环与仅目标节点行为对齐一维 API

有环、以及仅作为边目标出现的节点，分层结果中「哪些节点出现」MUST 与同图、同 `container` 下 `TopoSort` 一致（能排出的出现；环上无法排出的不出现）。层划分 MUST 仍遵循 Kahn 波次语义。

#### Scenario: 有环时仅返回能排出的层

- **WHEN** 图含环（例如 `a↔b` 且 `c→a`）并调用 `TopoSortLayers`（建议使用确定性 container 如 MinHeap）
- **THEN** 结果中仅出现能排出的节点（与 `TopoSort` 相同，例如仅 `c`），且这些节点按合法分层排列；环上无法排出的节点不出现

#### Scenario: 仅作为边目标的节点被包含

- **WHEN** 图为 `{"a":{"b"}}`（`b` 仅作为边目标）
- **THEN** 分层结果包含 `a` 与 `b`，且 `a` 所在层号严格小于 `b`（例如 `[["a"],["b"]]`）

### Requirement: DataByFormers 侧对称分层 API

`TopoSortDataByFormersLayers` SHALL 使用与 `TopoSortDataByFormers` 相同的泛型参数与构图规则（含忽略不在 `datas` 中的 former），对 key 图调用分层拓扑排序后映射回 `[][]Data`。

#### Scenario: 按 former 依赖分层

- **WHEN** 数据为 `a` 依赖 `b,c`，`b` 依赖 `c`，`c` 无依赖，并调用 `TopoSortDataByFormersLayers`（默认或指定 container）
- **THEN** 返回多层 `[][]Data`，按 key 观察时层间满足依赖先后（例如 `c` 在 `b` 之前层，`b` 在 `a` 之前层），且同一层内无 former 依赖边

#### Scenario: 未知 former 被忽略

- **WHEN** 某数据声明的 former 不在输入集合中（例如 `a` 的 former 为 `x`，另有 `b` 依赖 `a`）
- **THEN** 未知 former 不参与构图；结果中仍包含可排出的数据（与一维 `TopoSortDataByFormers` 对齐，例如包含 `a` 与 `b`）

#### Scenario: 数据侧有环

- **WHEN** 数据形成环（例如 `a↔b`）并另有独立节点 `c`
- **THEN** 分层结果仅包含能排出的数据（例如仅 `c`），与一维 API 对齐

### Requirement: 现有一维 API 不受影响

现有 `TopoSort` 与 `TopoSortDataByFormers` 的签名与行为 MUST 保持不变；既有 `stl` 包中相关单元测试 MUST 全部通过。

#### Scenario: 回归既有拓扑排序测试

- **WHEN** 完成本变更实现后运行 `go test ./stl/...`
- **THEN** 既有一维拓扑排序相关测试全部通过，且新增分层测试覆盖本 spec 中的场景
