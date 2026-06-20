package examples

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/fasionchan/goutils/libs/gonumplot"
)

//go:embed testdata/*.json
var testdata embed.FS

// LoadChartFixture 从 examples/testdata 加载 Chart JSON。
func LoadChartFixture(name string) (*gonumplot.Chart, error) {
	data, err := testdata.ReadFile("testdata/" + name)
	if err != nil {
		return nil, fmt.Errorf("read chart fixture %q: %w", name, err)
	}
	return gonumplot.LoadChartJSON(bytes.NewReader(data))
}
