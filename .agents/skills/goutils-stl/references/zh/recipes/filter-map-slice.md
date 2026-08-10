# 配方：过滤后再映射切片

**主主题：** slice · **符号：** `Filter`、`Map`、`Reduce`

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

**误用：** 不要期望切片版 `Map` 返回 `map[K]V`——应使用 `BuildMap` / `MappingByKey`。
