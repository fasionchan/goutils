# Recipe: build a map from a slice

**Primary:** map · **Symbols:** `MappingByKey`, `BuildMap`, `MapKeys`

```go
package main

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

type User struct {
	ID   string
	Name string
}

func main() {
	users := []User{{"u1", "Ada"}, {"u2", "Bob"}}
	byID := stl.MappingByKey(users, func(u User) string { return u.ID })
	fmt.Println(byID["u1"].Name) // Ada

	lengths := stl.BuildMap[map[string]int](users, func(u User) (string, int) {
		return u.ID, len(u.Name)
	})
	fmt.Println(stl.MapKeys(lengths))
}
```

**Misuse:** last-wins on duplicate keys for `MappingByKey` / `BuildMap` — dedupe first if needed (`UniqueByKeySet`).
