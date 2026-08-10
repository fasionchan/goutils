# 配方：将 iter.Seq 映射为切片

**主主题：** seq · **符号：** `DataSeq`、`MapSeqToSlice`、`ReadSeq`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

func main() {
	seq := stl.DataSeq(1, 2, 3)
	out := stl.MapSeqToSlice(seq, func(n int) int { return n * 10 })
	fmt.Println(out) // [10 20 30]

	again := stl.ReadSeq[[]int](stl.DataSeq(4, 5))
	fmt.Println(again) // [4 5]
}
```

**误用：** `iter.Seq` 单次遍历——不要对同一 seq 两次 `range`；需要时重新 `DataSeq` / `ReadSeq`。
