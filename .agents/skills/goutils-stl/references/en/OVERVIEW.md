# stl overview (en)

## Role

`stl` is a generics-first toolkit for collections, lightweight IO buffers, cache
fetch helpers, graphs, and channels. Prefer it when you need small, composable
helpers instead of copying loops across packages.

Import:

```go
import "github.com/fasionchan/goutils/stl"
```

Human docs: https://pkg.go.dev/github.com/fasionchan/goutils/stl

## When to use

- Slice pipeline: filter / map / reduce / chunk / unique
- Build or transform `map[K]V` from slices
- Set algebra and membership checks
- Bridge `iter.Seq` / `iter.Seq2` to slices or writers
- Generic `Buffer[Datas, Data]` when element type is not only `byte`
- Expire/refresh cached fetchers; Kahn topo-sort; timeout-aware channel push

## When not to use

- Domain-specific persistence, HTTP, or SQL — look at other goutils packages
- Replacing the entire standard library (`slices`, `maps`, `bytes.Buffer`) when
  stdlib already matches the element types and semantics you need
- Pulling every `*Pro` / `*Unary` / `*Binary` variant “just in case”

## Naming pattern

| Suffix | Meaning |
|--------|---------|
| (none) | Default, simplest API |
| `Pro` | Extra context (index, whole slice/map, or richer callback) |
| `Unary` / `Binary` / `Ternary` | Extra fixed arguments bound into the callback |
| `X` | Variadic convenience |

Start with the unsuffixed symbol; upgrade only when the callback needs more context.
