# goutils/stl index (en)

Progressive entry for `github.com/fasionchan/goutils/stl`.

## Start here

1. [OVERVIEW.md](OVERVIEW.md) — role, boundaries, discovery
2. Pick one capability (do not preload all):
   - [capabilities/slice.md](capabilities/slice.md) — **required hotspot**
   - [capabilities/map.md](capabilities/map.md) — **required hotspot**
   - [capabilities/set.md](capabilities/set.md)
   - [capabilities/seq.md](capabilities/seq.md)
   - [capabilities/buffer-io.md](capabilities/buffer-io.md)
   - [capabilities/cacher.md](capabilities/cacher.md)
   - [capabilities/graph.md](capabilities/graph.md)
   - [capabilities/chan.md](capabilities/chan.md)
3. Apply a recipe, then skim [antipatterns.md](antipatterns.md)

## Recipes

| Recipe | Primary area |
|--------|----------------|
| [recipes/filter-map-slice.md](recipes/filter-map-slice.md) | slice |
| [recipes/build-map-from-slice.md](recipes/build-map-from-slice.md) | map |
| [recipes/set-membership.md](recipes/set-membership.md) | set |
| [recipes/seq-to-slice.md](recipes/seq-to-slice.md) | seq |
| [recipes/buffer-generic-io.md](recipes/buffer-generic-io.md) | buffer/io |
| [recipes/topo-sort-by-formers.md](recipes/topo-sort-by-formers.md) | graph |

## Roadmap (other packages)

Later skills may cover `baseutils`, `std/*`, `queryutils`, etc. This skill is `stl` only.
