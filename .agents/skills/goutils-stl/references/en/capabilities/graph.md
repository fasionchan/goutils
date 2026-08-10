# Capability: graph

Source focus: `stl/graph.go`, containers in `stl/container.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Adjacency map | `Graph`, `GraphFromFormers` |
| Topo order | `Graph.TopoSort`, `Graph.TopoSortLayers` |
| Data-level sort | `TopoSortDataByFormers`, `TopoSortDataByFormersLayers` |
| Ready-set container | `NewStackAsContainer`, `NewMinHeapAsContainer`, `NewMaxHeapAsContainer` |

Pass `nil` container factory for default stack ordering; pass heap factories when
tie-breaking by key order matters.
