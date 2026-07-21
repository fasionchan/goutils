package browser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryEncodeEnbed(t *testing.T) {
	query, err := queryEncoder.Encode(struct {
		Clip      *Viewport `query:"clip"`
		*Viewport `query:",inline"`
	}{
		Clip: &Viewport{
			X:      10,
			Y:      20,
			Width:  30,
			Height: 40,
		},
		Viewport: &Viewport{
			X:      100,
			Y:      200,
			Width:  300,
			Height: 400,
		},
	})
	if err != nil {
		t.Fatalf("Failed to encode query: %v", err)
	}

	fmt.Println(query.Encode())

	require.Equal(t, query.Get("clip.x"), "10")
	require.Equal(t, query.Get("clip.y"), "20")
	require.Equal(t, query.Get("clip.width"), "30")
	require.Equal(t, query.Get("clip.height"), "40")
	require.Equal(t, query.Get("clip.scale"), "1")

	require.Equal(t, query.Get("x"), "100")
	require.Equal(t, query.Get("y"), "200")
	require.Equal(t, query.Get("width"), "300")
	require.Equal(t, query.Get("height"), "400")
	require.Equal(t, query.Get("scale"), "2")
}
