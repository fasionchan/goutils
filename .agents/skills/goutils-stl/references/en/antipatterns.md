# Antipatterns (en)

1. **Load both `en/` and `zh/`** — wastes tokens; pick one language tree.
2. **Paste the whole skill into prompts** — use progressive disclosure instead.
3. **Reimplement `Filter`/`Map`/`Contain`** — check capabilities first.
4. **Confuse slice `Map` with map package helpers** — `stl.Map` transforms slices;
   use `BuildMap` / `MappingByKey` / `MapMap*` for `map[K]V`.
5. **Ignore nil / empty semantics** — many helpers treat nil/empty slices as empty
   results; verify with tests when zero values matter.
6. **Use blocking `Chan.Push` in request paths** — prefer `PushPro` / `ChanPipe`
   with context when cancellation matters.
7. **Treat this skill as godoc** — confirm signatures with `go doc` before shipping.
