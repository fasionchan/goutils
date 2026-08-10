# 配方：非 byte 切片的泛型 Buffer

**主主题：** buffer/io · **符号：** `NewBuffer`、`Write`、`Datas`

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

**误用：** 普通 `[]byte` 文本/二进制通常 `bytes.Buffer` 即可；元素不是 `byte`、或已统一用 `stl.Writer` 时再用 `stl.Buffer`。
