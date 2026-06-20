package gonumplot

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderChartPNG(t *testing.T) {
	chart := &Chart{
		Title:  "Test",
		XLabel: "X",
		YLabel: "Y",
		Series: []Series{
			{
				Name: "line1",
				Points: []Point{
					{X: 1, Y: 1},
					{X: 2, Y: 3},
					{X: 3, Y: 2},
				},
			},
		},
	}

	pngData, err := RenderChartPNG(chart)
	require.NoError(t, err)
	require.NotEmpty(t, pngData)
	require.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47}, pngData[:4])
}

func TestRenderChartMultiSeries(t *testing.T) {
	chart := &Chart{
		Title:  "Multi",
		XLabel: "X",
		YLabel: "Y",
		Series: []Series{
			{Name: "a", Points: []Point{{X: 1, Y: 1}, {X: 2, Y: 2}}},
			{Name: "b", Points: []Point{{X: 1, Y: 2}, {X: 2, Y: 1}}},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, RenderChart(chart, &buf))
	require.NotEmpty(t, buf.Bytes())
}

func TestRenderChartEmpty(t *testing.T) {
	_, err := RenderChartPNG(&Chart{})
	require.ErrorIs(t, err, errEmptyChart)

	_, err = RenderChartPNG(&Chart{Series: []Series{{Name: "empty"}}})
	require.ErrorIs(t, err, errEmptyChart)
}
