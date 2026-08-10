# 能力：slice

源码重点：`stl/slice.go`、`stl/slice_pro.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 任意/全部谓词 | `AnyMatch`、`AllMatch` |
| 过滤 | `Filter`、`FilterByKey`、`FilterByKeys` |
| 切片映射 | `Map`、`MapUntilError`、`MapAndConcat` |
| 归约 | `Reduce` |
| 查找 | `FindFirst`、`FindFirstOrZero`、`FindLast` |
| 成员判断 | `Contain`、`ContainAll`、`ContainAny`、`IndexOf` |
| 分块 | `Divide` |
| 拼接 | `ConcatSlices`、`JoinSlices` |
| 去重 | `UniqueBySet`、`StableUniqueBySet`、`UniqueByKeySet` |
| 安全取下标 | `Index`、`FirstOneOrZero`、`LastOneOrZero` |
| 排序相关 | `Sort` 及比较助手（按需） |

回调需要下标或额外参数时再用 `*Pro` / `*Unary`。

## 说明

- 这里的 `Map` 变换的是**切片元素**，不是 `map[K]V`。
- 不关心顺序用 `UniqueBySet`；要保留首次出现顺序用 `StableUniqueBySet`。
