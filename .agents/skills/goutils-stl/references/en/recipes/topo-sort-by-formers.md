# Recipe: topological sort by former keys

**Primary:** graph · **Symbols:** `TopoSortDataByFormers`, `Map`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

type Node struct {
	Key     string
	Formers []string
}

func main() {
	nodes := []Node{
		{Key: "a", Formers: []string{"b", "c"}},
		{Key: "b", Formers: []string{"c"}},
		{Key: "c"},
	}
	ordered := stl.TopoSortDataByFormers(
		nodes,
		func(n Node) string { return n.Key },
		func(n Node) []string { return n.Formers },
		nil, // default stack container
	)
	fmt.Println(stl.Map(ordered, func(n Node) string { return n.Key })) // [c b a]
}
```

**Misuse:** cycles yield incomplete orders — validate graph assumptions; pass `NewMinHeapAsContainer[string]` when deterministic key order is required.
