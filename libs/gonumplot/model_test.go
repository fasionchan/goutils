package gonumplot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChartJSONRoundTrip(t *testing.T) {
	original := &Chart{
		Title:  "CPU Usage",
		XLabel: "Time",
		YLabel: "%",
		Series: []Series{
			{
				Name: "cpu0",
				Points: []Point{
					{X: 1718000000, Y: 45.2},
					{X: 1718000060, Y: 50.1},
				},
			},
		},
	}

	data, err := MarshalChartJSON(original)
	require.NoError(t, err)

	loaded, err := LoadChartJSON(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, original.Title, loaded.Title)
	require.Equal(t, original.XLabel, loaded.XLabel)
	require.Equal(t, original.YLabel, loaded.YLabel)
	require.Len(t, loaded.Series, 1)
	require.Equal(t, original.Series[0].Name, loaded.Series[0].Name)
	require.Equal(t, original.Series[0].Points, loaded.Series[0].Points)
}

func TestLoadChartJSONInvalid(t *testing.T) {
	_, err := LoadChartJSON(strings.NewReader("{invalid"))
	require.Error(t, err)
}

func TestMarshalChartJSON(t *testing.T) {
	chart := &Chart{Title: "test", Series: []Series{{Name: "a", Points: []Point{{X: 1, Y: 2}}}}}
	data, err := MarshalChartJSON(chart)
	require.NoError(t, err)
	require.True(t, json.Valid(data))
}
