# Capability: cacher

Source focus: `stl/cacher.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Construct | `NewCachedDataFetcher`, `NewCachedDataFetcherLite` |
| Fetch | `Fetch`, `FetchWithExpires` (methods on `CachedDataFetcher`) |
| Derive | `NewCachedDataFetcherFromAnother` |
| Callbacks | `NewCachedDataFetcherCallback`, callback registration methods |
| Adapter | `GetCachedDataFetchLite` |

Use for in-process cached fetches with expiry/refresh — not a distributed cache.
