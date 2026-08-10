# 配方：集合成员与运算

**主主题：** set · **符号：** `NewSet`、`Contain`、`Intersection`

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
	fmt.Println(allowed.Intersection(requested).Slice()) // [write]（顺序不保证）
}
```

**误用：** `Set` 不保证顺序——仅在不关心顺序或自行排序时使用 `Slice()`。
