# 反模式（zh）

1. **同时加载 `en/` 与 `zh/`** — 浪费 token；只选一棵语言树。
2. **把整个 Skill 粘进提示词** — 应走渐进加载。
3. **手写 `Filter`/`Map`/`Contain`** — 先查 capabilities。
4. **混淆切片 `Map` 与 map 助手** — `stl.Map` 变换的是切片；
   `map[K]V` 用 `BuildMap` / `MappingByKey` / `MapMap*`。
5. **忽略 nil/空切片语义** — 多数助手把 nil/空当作空结果；涉及零值时补测试。
6. **在请求路径用阻塞 `Chan.Push`** — 需要取消时用 `PushPro` / `ChanPipe`。
7. **把本 Skill 当 godoc** — 上线前用 `go doc` 核对签名。
