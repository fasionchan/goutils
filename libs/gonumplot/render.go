package gonumplot

import (
	"bytes"
	"fmt"
	"io"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	defaultWidth  = 6 * vg.Inch
	defaultHeight = 4 * vg.Inch
)

// RenderChart 将 Chart 渲染为 PNG 并写入 w。
func RenderChart(chart *Chart, w io.Writer) error {
	if err := validateChart(chart); err != nil {
		return err
	}

	p, err := newPlot(chart)
	if err != nil {
		return err
	}

	writerTo, err := p.WriterTo(defaultWidth, defaultHeight, "png")
	if err != nil {
		return fmt.Errorf("create png writer: %w", err)
	}
	if _, err := writerTo.WriteTo(w); err != nil {
		return fmt.Errorf("write png: %w", err)
	}
	return nil
}

// RenderChartPNG 将 Chart 渲染为 PNG 字节切片。
func RenderChartPNG(chart *Chart) ([]byte, error) {
	var buf bytes.Buffer
	if err := RenderChart(chart, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func validateChart(chart *Chart) error {
	if chart == nil || len(chart.Series) == 0 {
		return errEmptyChart
	}
	hasPoints := false
	for _, series := range chart.Series {
		if len(series.Points) > 0 {
			hasPoints = true
			break
		}
	}
	if !hasPoints {
		return errEmptyChart
	}
	return nil
}

func newPlot(chart *Chart) (*plot.Plot, error) {
	p := plot.New()
	p.Title.Text = chart.Title
	p.X.Label.Text = chart.XLabel
	p.Y.Label.Text = chart.YLabel

	loc, err := chartLocation(chart)
	if err != nil {
		return nil, err
	}
	p.X.Tick.Marker = newDatetimeTicker(loc)

	args := make([]any, 0, len(chart.Series)*2)
	for _, series := range chart.Series {
		if len(series.Points) == 0 {
			continue
		}
		pts := make(plotter.XYs, len(series.Points))
		for i, point := range series.Points {
			pts[i].X = point.X
			pts[i].Y = point.Y
		}
		args = append(args, series.Name, pts)
	}
	if len(args) == 0 {
		return nil, errEmptyChart
	}
	if err := plotutil.AddLines(p, args...); err != nil {
		return nil, fmt.Errorf("add lines: %w", err)
	}
	return p, nil
}
