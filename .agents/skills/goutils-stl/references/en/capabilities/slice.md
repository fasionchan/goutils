# Capability: slice

Source focus: `stl/slice.go`, `stl/slice_pro.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Any / all predicate | `AnyMatch`, `AllMatch` |
| Filter | `Filter`, `FilterByKey`, `FilterByKeys` |
| Map slice → slice | `Map`, `MapUntilError`, `MapAndConcat` |
| Reduce | `Reduce` |
| Find | `FindFirst`, `FindFirstOrZero`, `FindLast` |
| Membership | `Contain`, `ContainAll`, `ContainAny`, `IndexOf` |
| Chunk | `Divide` |
| Concat | `ConcatSlices`, `JoinSlices` |
| Unique | `UniqueBySet`, `StableUniqueBySet`, `UniqueByKeySet` |
| Index safely | `Index`, `FirstOneOrZero`, `LastOneOrZero` |
| Sort helpers | `Sort`, compare helpers as needed |

`*Pro` / `*Unary` variants exist when the callback needs index or extra args —
start without them.

## Notes

- `Map` here maps **slice elements**, not `map[K]V`.
- Prefer `UniqueBySet` when order does not matter; `StableUniqueBySet` keeps first-seen order.
