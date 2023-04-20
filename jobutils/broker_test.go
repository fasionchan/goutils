/*
 * Author: fasion
 * Created time: 2023-04-20 17:23:45
 * Last Modified by: fasion
 * Last Modified time: 2023-04-20 17:35:06
 */

package jobutils

import (
	"fmt"
	"testing"
	"time"
)

func TestTopicBroker(t *testing.T) {
	broker := NewTopicBroker()
	broker.Subscribe("test", func() {
		fmt.Println("callAt:", time.Now())
	}).
		WithAsync(true).
		WithConcurrentcy(1).
		WithInterval(time.Second).
		WithMerge(true).
		Done()

	for i := 0; i < 100; i++ {
		broker.Publish("test")
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println("waiting....")

	time.Sleep(time.Minute)
}
