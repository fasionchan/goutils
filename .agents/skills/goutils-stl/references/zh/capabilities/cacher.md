# 能力：cacher

源码重点：`stl/cacher.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 构造 | `NewCachedDataFetcher`、`NewCachedDataFetcherLite` |
| 拉取 | `Fetch`、`FetchWithExpires`（`CachedDataFetcher` 方法） |
| 派生 | `NewCachedDataFetcherFromAnother` |
| 回调 | `NewCachedDataFetcherCallback` 及注册方法 |
| 适配 | `GetCachedDataFetchLite` |

用于进程内带过期/刷新的缓存拉取，不是分布式缓存。
