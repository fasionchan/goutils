# Recipe: generic Buffer for non-byte slices

**Primary:** buffer/io · **Symbols:** `NewBuffer`, `Write`, `Datas`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

func main() {
	buf := stl.NewBuffer[[]int]()
	_, _ = buf.Write([]int{1, 2})
	_, _ = buf.Write([]int{3})
	fmt.Println(buf.Datas()) // [1 2 3]
}
```

**Misuse:** for plain `[]byte` text/binary, `bytes.Buffer` is often enough; use `stl.Buffer` when the element type is not `byte` or you already standardize on `stl.Writer`.
