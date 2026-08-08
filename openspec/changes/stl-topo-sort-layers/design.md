## Context

`stl/graph.go` 已实现 Kahn 风格一维拓扑排序：

- `Graph[T].TopoSort(container)` → `[]T`
- `TopoSortDataByFormers(...)` → 构图后调用 `TopoSort`，再映射回 `[]Data`

入度统计、默认 `container`（`nil` → `NewStackAsContainer`）、仅边目标节点、未知 former 忽略、有环时只返回能排出的前缀等行为已由 `stl/graph_test.go` 锁定。

需求分析（FAS-2）已确认：「分层」= Kahn 每轮入度为 0 的一批；保留现有 API；需要 DataByFormers 对称版本；环/未知 former 完全对齐；层内顺序沿用 container。实现由 fasionchan 跟进。

## Goals / Non-Goals

**Goals:**

- 新增 `TopoSortLayers` / `TopoSortDataByFormersLayers`，返回 `[][]T` / `[][]Data`
- 层语义：同一层内节点彼此无依赖边，可并行；层间满足拓扑先后（对任意边 `u→v`，`u` 的层号严格小于 `v`）
- 与现有一维 API 在入度初始化、默认 container、环、未知 former、仅目标节点等边界上对齐
- 单测覆盖 AC-1～AC-7；既有一维测试全部保持通过

**Non-Goals:**

- 不改现有 `TopoSort` / `TopoSortDataByFormers` 签名或语义
- 不实现「最长路径深度分层」等非 Kahn 语义
- 不引入可视化、持久化或新外部依赖
- 不保证 flatten(`layers`) 与同参数 `TopoSort` 的全局序列完全一致（一维会即时把新零入度节点混入同一队列；分层按波次排出，全局交错不同属预期）

## Decisions

### 1. API 形态：新增方法，不改返回类型

**选择**:

```go
func (g Graph[T]) TopoSortLayers(container func(capacity int) Container[T]) [][]T

func TopoSortDataByFormersLayers[
	Datas ~[]Data,
	Keys ~[]Key,
	Data any,
	Key comparable,
](datas Datas, getKey func(Data) Key, getFormerKeys func(Data) Keys, container func(capacity int) Container[Key]) [][]Data
```

**理由**: 兼容现有调用方；与已确认决策「保留 TopoSort」一致；Data 侧与现有对称，满足 AC-7。

**备选**: 修改 `TopoSort` 返回 `[][]T`（**BREAKING**，已否决）；仅 Graph 侧新增（数据侧主要用法缺失，已否决）。

### 2. 层算法：整波排空再推进（禁止按 size 混推）

**选择**: 每一轮先把当前 container 中**全部**节点 Pop 到本层切片（层内顺序即连续 Pop 顺序），再统一处理这些节点的出边、把新入度为 0 的节点 Push 进 container，作为下一层候选。

**理由**: 默认 container 是栈（LIFO）。若采用「记录本轮 size、边处理时立刻 Push」的经典队列写法，新节点会被先于本层剩余节点 Pop 出，破坏「同一层」与层内 container 语义。整波排空对 Stack / Heap / 任意 Container 均正确。

**备选**: 双缓冲（current/next 两个 container）—— 等价但分配更多；本设计用「先排空再推进」即可。

伪代码：

```text
build inDegree  // 与 TopoSort 相同：先初始化 g 的 key，再对边目标 ++
q := container(初始零入度节点)
layers := []
for !q.IsEmpty():
  layer := []
  for !q.IsEmpty():
    layer = append(layer, q.Pop())
  for node in layer:
    for neighbor in g[node]:
      inDegree[neighbor]--
      if inDegree[neighbor] == 0:
        q.Push(neighbor)
  layers = append(layers, layer)
return layers
```

### 3. 入度与图构建：复用现有逻辑，避免漂移

**选择**:

- `TopoSortLayers` 的入度初始化与边遍历 MUST 与 `TopoSort` 一致（含仅出现在边目标的节点）。
- `TopoSortDataByFormersLayers` 的构图 MUST 与 `TopoSortDataByFormers` 一致（未知 former 跳过；`graph[key]` 保证存在）。

**理由**: AC-5 要求「完全对齐」；复制同一构图路径可避免行为漂移。允许在同文件内抽取私有辅助函数（如 `buildInDegree` / 构图），但抽取不得改变现有一维公开 API 行为（既有测试必须仍通过）。

**备选**: Layers 内部直接调用 `TopoSort` 再二次分层 —— 不可行，一维结果丢失波次边界。

### 4. 空结果惯例

**选择**: 空图返回长度为 0 的 `[][]T`（可用 `make([][]T, 0)`）；单测用 `assert.Empty`，不强制区分 `nil` 与非 nil 空切片。

**理由**: 与现有 `TopoSort` 使用 `make([]T, 0)` 及 `assert.Empty` 一致；AC-4 已允许二选一。

### 5. Data 侧映射

**选择**: 对 `TopoSortLayers` 得到的每一层 `[]Key`，用现有 `MapValuesByKeys`（或等价逻辑）映射为 `[]Data`，组装为 `[][]Data`。

**理由**: 与一维 `MapValuesByKeys(mapping, graph.TopoSort(...)...)` 对称，保证 key→data 一致性。

### 6. 代码落点

**选择**: 实现放在 `stl/graph.go`（或同包极小拆分如 `graph_layers.go`，若希望文件更短）；测试追加到 `stl/graph_test.go`。

**理由**: 与现有拓扑排序同文件同包，调用方 import 路径不变；最小变更。

## Risks / Trade-offs

- **[Risk] 实现时误用「按 size 边推」导致栈下层语义错误** → 设计强制整波排空；单测用默认栈 + 多源图断言层边界
- **[Risk] 抽取公共入度逻辑时意外改坏 `TopoSort`** → 先跑既有 `graph_test.go`；抽取后 diff 行为为零
- **[Risk] 调用方误以为 flatten(layers) == TopoSort(...)** → 在 Go doc 中写明分层为 Kahn 波次，全局交错可能不同于一维
- **[Trade-off] 不保证层内稳定排序（无 container 时依赖 map 迭代）** → 与现有 `TopoSort(nil)` 一致；需要确定性时传入 Min/MaxHeap container

## Migration Plan

纯新增 API，无需迁移或回滚脚本。若需回退：删除新函数与对应测试即可，不影响现有调用方。

## Open Questions

无阻塞项。以下为实现期约定（已可按此编码）：

- flatten 与一维顺序不一致：接受，不单测强制相等。
- 是否抽取私有辅助：实现者自定；以既有一维测试全绿为约束。
