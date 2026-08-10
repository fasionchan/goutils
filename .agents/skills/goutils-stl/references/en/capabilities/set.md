# Capability: set

Source focus: `stl/set.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Construct | `NewSet`, `NewEmptySet` |
| Membership | `Contain`, `ContainAll`, `ContainAny` |
| Mutate | `Push` / `PushX`, `Add` / `AddX`, `Pop`, `Merge`, `Purge` |
| Algebra | `Union`, `Intersection`, `Difference`, `SymmetricDifference` |
| Export | `Slice`, `Dup`, `Equal`, `Len`, `Empty` |

Fewer symbols than slice/map by design — set is a thin `map[T]struct{}` wrapper.
