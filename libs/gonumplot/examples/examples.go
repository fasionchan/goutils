package examples

import (
	"io"

	"github.com/fasionchan/goutils/libs/gonumplot"
)

// RenderCPU 加载 CPU 监控 Chart fixture 并渲染为 PNG。
func RenderCPU(w io.Writer) error {
	chart, err := LoadChartFixture("cpu.json")
	if err != nil {
		return err
	}
	return gonumplot.RenderChart(chart, w)
}

// RenderMemory 加载内存监控 Chart fixture 并渲染为 PNG。
func RenderMemory(w io.Writer) error {
	chart, err := LoadChartFixture("memory.json")
	if err != nil {
		return err
	}
	return gonumplot.RenderChart(chart, w)
}

// RenderFilesystem 加载文件系统空间监控 Chart fixture 并渲染为 PNG。
func RenderFilesystem(w io.Writer) error {
	chart, err := LoadChartFixture("filesystem.json")
	if err != nil {
		return err
	}
	return gonumplot.RenderChart(chart, w)
}

// RenderDiskIO 加载磁盘 I/O 监控 Chart fixture 并渲染为 PNG。
func RenderDiskIO(w io.Writer) error {
	chart, err := LoadChartFixture("diskio.json")
	if err != nil {
		return err
	}
	return gonumplot.RenderChart(chart, w)
}

// RenderNetwork 加载网络流量监控 Chart fixture 并渲染为 PNG。
func RenderNetwork(w io.Writer) error {
	chart, err := LoadChartFixture("network.json")
	if err != nil {
		return err
	}
	return gonumplot.RenderChart(chart, w)
}
