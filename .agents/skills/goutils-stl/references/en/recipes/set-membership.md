# Recipe: set membership and algebra

**Primary:** set · **Symbols:** `NewSet`, `Contain`, `Intersection`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

func main() {
	allowed := stl.NewSet("read", "write")
	fmt.Println(allowed.Contain("write")) // true

	requested := stl.NewSet("write", "admin")
	fmt.Println(allowed.Intersection(requested).Slice()) // [write] (order not guaranteed)
}
```

**Misuse:** do not use `Set` for ordered collections — convert with `Slice()` only when order is irrelevant or you sort afterward.
