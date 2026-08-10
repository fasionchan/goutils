# 配方：按前置键拓扑排序

**主主题：** graph · **符号：** `TopoSortDataByFormers`、`Map`

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
		nil, // 默认栈容器
	)
	fmt.Println(stl.Map(ordered, func(n Node) string { return n.Key })) // [c b a]
}
```

**误用：** 有环时结果可能不完整——先确认图假设；需要按 key 稳定决胜时传入 `NewMinHeapAsContainer[string]`。
