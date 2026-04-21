package queryutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

}

func TestMultiSetinHandler(t *testing.T) {
	type Data = struct{}
	type Datas = []*Data

	var handler1Ok bool
	handler1 := func(ctx context.Context, datas Datas) error {
		handler1Ok = true
		return nil
	}

	var handler2Ok bool
	handler2 := func(ctx context.Context, datas Datas) error {
		handler2Ok = true
		return nil
	}

	handler := MultiSetinHandler(handler1, handler2)
	err := handler(context.Background(), []*Data{&Data{}})
	if err != nil {
		t.Fatalf("MultiSetinHandler failed: %v", err)
	}

	assert.True(t, handler1Ok)
	assert.True(t, handler2Ok)
}