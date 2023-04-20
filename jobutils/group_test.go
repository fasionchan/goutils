/*
 * Author: fasion
 * Created time: 2023-03-15 14:38:32
 * Last Modified by: fasion
 * Last Modified time: 2023-03-15 17:25:06
 */

package jobutils

import (
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	tokens := NewJobTokens(10)

	if tokens.Totals() != 10 {
		t.Fatal("totals")
	}
	if tokens.Useds() != 0 {
		t.Fatal("useds")
	}
	if tokens.Lefts() != 10 {
		t.Fatal("lefts")
	}

	tokens.Acquire(nil, -1)
	if tokens.Totals() != 10 {
		t.Fatal("totals")
	}
	if tokens.Useds() != 1 {
		t.Fatal("useds")
	}
	if tokens.Lefts() != 9 {
		t.Fatal("lefts")
	}

	for i := 0; i < 9; i++ {
		tokens.Acquire(nil, -1)
	}

	if tokens.Acquire(nil, 0) {
		t.Fatal("Acquire ok")
	}

	startTime := time.Now()
	if tokens.Acquire(nil, time.Second) {
		t.Fatal("Acquire ok")
	}
	delta := time.Now().Sub(startTime) - time.Second
	if delta < 0 {
		delta = -delta
	}
	if delta > time.Millisecond {
		t.Fatal("Acquire wait error")
	}
}
