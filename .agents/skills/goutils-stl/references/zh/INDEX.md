# goutils/stl 索引（zh）

`github.com/fasionchan/goutils/stl` 的渐进式入口。

## 从这里开始

1. [OVERVIEW.md](OVERVIEW.md) — 职责、边界、发现方式
2. 只打开一个能力主题（不要预加载全部）：
   - [capabilities/slice.md](capabilities/slice.md) — **必读热点**
   - [capabilities/map.md](capabilities/map.md) — **必读热点**
   - [capabilities/set.md](capabilities/set.md)
   - [capabilities/seq.md](capabilities/seq.md)
   - [capabilities/buffer-io.md](capabilities/buffer-io.md)
   - [capabilities/cacher.md](capabilities/cacher.md)
   - [capabilities/graph.md](capabilities/graph.md)
   - [capabilities/chan.md](capabilities/chan.md)
3. 套用一条 recipe，再扫一眼 [antipatterns.md](antipatterns.md)

## 配方

| 配方 | 主主题 |
|------|--------|
| [recipes/filter-map-slice.md](recipes/filter-map-slice.md) | slice |
| [recipes/build-map-from-slice.md](recipes/build-map-from-slice.md) | map |
| [recipes/set-membership.md](recipes/set-membership.md) | set |
| [recipes/seq-to-slice.md](recipes/seq-to-slice.md) | seq |
| [recipes/buffer-generic-io.md](recipes/buffer-generic-io.md) | buffer/io |
| [recipes/topo-sort-by-formers.md](recipes/topo-sort-by-formers.md) | graph |

## 路线图（其他包）

后续可能覆盖 `baseutils`、`std/*`、`queryutils` 等。本 Skill 仅针对 `stl`。
