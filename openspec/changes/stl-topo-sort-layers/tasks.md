## 1. Graph 侧实现

- [x] 1.1 在 `stl/graph.go` 实现 `Graph[T].TopoSortLayers`：入度初始化与 `TopoSort` 对齐
- [x] 1.2 按「整波排空再推进」实现 Kahn 分层；`container == nil` 时默认 `NewStackAsContainer`
- [x] 1.3 为空图返回空分层结果；为公开方法补充简要 Go doc（说明 Kahn 波次语义，及 flatten 不必等于一维顺序）

## 2. DataByFormers 侧实现

- [x] 2.1 实现 `TopoSortDataByFormersLayers`，构图规则与 `TopoSortDataByFormers` 一致（未知 former 跳过）
- [x] 2.2 将各层 `[]Key` 映射为 `[]Data`（复用 `MapValuesByKeys` 或等价逻辑）

## 3. 可选重构（不得改变一维行为）

- [x] 3.1 若需要，抽取私有入度/构图辅助函数供一维与分层共用
- [x] 3.2 抽取后立即运行既有 `TestTopoSort*` / `TestGraphTopoSort*` 确认无回归

## 4. 单元测试

- [x] 4.1 空图与单节点：`TopoSortLayers` 边界（AC-4）
- [x] 4.2 无环多源图：断言层间拓扑、层内无边；MinHeap 下 `a,b → c` 为 `[["a","b"],["c"]]`（AC-1/AC-2/AC-6）
- [x] 4.3 有环与仅目标节点：与同条件 `TopoSort` 的节点集合对齐（AC-5）
- [x] 4.4 `TopoSortDataByFormersLayers`：依赖链、未知 former、有环（AC-5/AC-7）
- [x] 4.5 `container == nil` 时结果合法（层约束成立即可）
- [x] 4.6 运行 `go test ./stl/...` 确保新旧测试全部通过（AC-3）

## 5. 收尾

- [x] 5.1 确认未修改现有 `TopoSort` / `TopoSortDataByFormers` 公开签名
- [ ] 5.2 按需开 PR（标题或 body 含 `FAS-2` / `Closes FAS-2`），便于 issue 关联 — 分支已推送 `fas-2-stl-topo-sort-layers`，本机缺 `gh` 登录未能自动建 PR
