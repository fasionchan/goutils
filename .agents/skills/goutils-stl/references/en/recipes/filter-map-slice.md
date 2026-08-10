# Recipe: filter then map a slice

**Primary:** slice · **Symbols:** `Filter`, `Map`, `Reduce`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

func main() {
	nums := []int{0, 1, 2, 3, 4, 5}
	evens := stl.Filter(nums, func(n int) bool { return n%2 == 0 })
	doubled := stl.Map(evens, func(n int) int { return n * 2 })
	sum := stl.Reduce(doubled, func(acc, n int) int { return acc + n }, 0)
	fmt.Println(evens, doubled, sum) // [0 2 4] [0 4 8] 12
}
```

**Misuse:** do not call slice `Map` expecting `map[K]V` output — use `BuildMap` / `MappingByKey`.
