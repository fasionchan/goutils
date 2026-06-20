package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadChartFixture(t *testing.T) {
	chart, err := LoadChartFixture("cpu.json")
	require.NoError(t, err)
	require.NotEmpty(t, chart.Title)
	require.NotEmpty(t, chart.Series)
}

func TestLoadChartFixtureAll(t *testing.T) {
	fixtures := []string{"cpu.json", "memory.json", "filesystem.json", "diskio.json", "network.json"}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			chart, err := LoadChartFixture(name)
			require.NoError(t, err)
			require.NotEmpty(t, chart.Series)
		})
	}
}
