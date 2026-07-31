package timex

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var i int
	for date := range PreviousDaySeq(time.Now()) {
		fmt.Println(date)
		i ++
		if i == 10 {
			break
		}
	}
}