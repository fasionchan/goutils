# 能力：set

源码重点：`stl/set.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 构造 | `NewSet`、`NewEmptySet` |
| 成员 | `Contain`、`ContainAll`、`ContainAny` |
| 变更 | `Push` / `PushX`、`Add` / `AddX`、`Pop`、`Merge`、`Purge` |
| 集合运算 | `Union`、`Intersection`、`Difference`、`SymmetricDifference` |
| 导出 | `Slice`、`Dup`、`Equal`、`Len`、`Empty` |

相对 slice/map 符号更少——`Set` 是薄的 `map[T]struct{}` 包装。
