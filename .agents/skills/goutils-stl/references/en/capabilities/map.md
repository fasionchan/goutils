# Capability: map

Source focus: `stl/map.go`, `stl/sync_map.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Typed map wrapper | `Mapping`, `NewMapping`, `NewMappingWithCap` |
| Slice → map | `BuildMap`, `BuildMapPro`, `MappingByKey`, `MappingByKeys` |
| Keys / values | `MapKeys`, `MapValues`, `MapValuesByKeys` |
| Filter map | `FilterMap`, `FilterMapByKey`, `FilterMapByValue` |
| Transform map | `MapMap`, `MapMapLite`, `MapMapToSlice` |
| Merge | `ConcatMap`, `ConcatMaps`, `ConcatMapInplace` |
| Subset / purge | `SubMapByKeys`, `BatchDeleteMap`, `PurgeMapKeys` |
| Lazy fill | `CacheMapValue`, `CacheMapValueWithInitializer`, `LoadOrCreate` |
| Counting | `Counter`, `Increase` / `Decrease` |
| Concurrent | `SyncMap`, `NewSyncMap`, `NewSyncMapPro` |

This topic is dense; recipes show the common slice→map path first.
