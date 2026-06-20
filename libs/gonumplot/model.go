package gonumplot

import (
	"encoding/json"
	"fmt"
	"io"
)

// Chart 表示多序列折线图的数据，支持 JSON 序列化。
type Chart struct {
	Title    string   `json:"title"`
	XLabel   string   `json:"xLabel"`
	YLabel   string   `json:"yLabel"`
	TimeZone string   `json:"timeZone,omitempty"` // IANA 时区，如 Asia/Shanghai；空为 UTC
	Series   []Series `json:"series"`
}

// Series 表示一条时序折线。
type Series struct {
	Name   string  `json:"name"`
	Points []Point `json:"points"`
}

// Point 表示折线上的一个点；X 通常为 Unix 时间戳（秒）。
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// LoadChartJSON 从 JSON 读取 Chart。
func LoadChartJSON(r io.Reader) (*Chart, error) {
	var chart Chart
	if err := json.NewDecoder(r).Decode(&chart); err != nil {
		return nil, fmt.Errorf("decode chart json: %w", err)
	}
	return &chart, nil
}

// MarshalChartJSON 将 Chart 序列化为 JSON。
func MarshalChartJSON(chart *Chart) ([]byte, error) {
	return json.MarshalIndent(chart, "", "  ")
}
