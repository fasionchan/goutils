# 能力：map

源码重点：`stl/map.go`、`stl/sync_map.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 带方法的 map 包装 | `Mapping`、`NewMapping`、`NewMappingWithCap` |
| 切片 → map | `BuildMap`、`BuildMapPro`、`MappingByKey`、`MappingByKeys` |
| 键 / 值列表 | `MapKeys`、`MapValues`、`MapValuesByKeys` |
| 过滤 map | `FilterMap`、`FilterMapByKey`、`FilterMapByValue` |
| 变换 map | `MapMap`、`MapMapLite`、`MapMapToSlice` |
| 合并 | `ConcatMap`、`ConcatMaps`、`ConcatMapInplace` |
| 子集 / 删除 | `SubMapByKeys`、`BatchDeleteMap`、`PurgeMapKeys` |
| 惰性填值 | `CacheMapValue`、`CacheMapValueWithInitializer`、`LoadOrCreate` |
| 计数 | `Counter`、`Increase` / `Decrease` |
| 并发 | `SyncMap`、`NewSyncMap`、`NewSyncMapPro` |

本主题符号较多；配方优先展示切片 → map 主路径。
