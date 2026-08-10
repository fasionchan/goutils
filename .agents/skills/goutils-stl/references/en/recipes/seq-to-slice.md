# Recipe: map an iter.Seq into a slice

**Primary:** seq · **Symbols:** `DataSeq`, `MapSeqToSlice`, `ReadSeq`

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

**Misuse:** `iter.Seq` is single-pass — do not range the same seq twice; rematerialize with `DataSeq` / `ReadSeq` as needed.
