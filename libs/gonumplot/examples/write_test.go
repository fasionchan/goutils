package examples

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// exampleOutputs 定义示例名称与输出文件名、渲染函数的对应关系。
var exampleOutputs = []struct {
	name   string
	file   string
	render func(io.Writer) error
}{
	{"cpu", "cpu.png", RenderCPU},
	{"memory", "memory.png", RenderMemory},
	{"filesystem", "filesystem.png", RenderFilesystem},
	{"diskio", "diskio.png", RenderDiskIO},
	{"network", "network.png", RenderNetwork},
}

// TestWriteExampleCharts 将各监控示例渲染为 PNG 文件，写入 examples/output/ 目录。
// 运行: go test ./libs/gonumplot/examples/ -run TestWriteExampleCharts -v
func TestWriteExampleCharts(t *testing.T) {
	outDir := filepath.Join("output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}

	for _, ex := range exampleOutputs {
		t.Run(ex.name, func(t *testing.T) {
			path := filepath.Join(outDir, ex.file)
			f, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			if err := ex.render(f); err != nil {
				t.Fatal(err)
			}
			t.Logf("written: %s", path)
		})
	}
}

// 以下测试可单独运行，便于只生成某一类图表：
// go test ./libs/gonumplot/examples/ -run TestWriteCPUChart -v

func TestWriteCPUChart(t *testing.T)        { writeChart(t, "cpu.png", RenderCPU) }
func TestWriteMemoryChart(t *testing.T)     { writeChart(t, "memory.png", RenderMemory) }
func TestWriteFilesystemChart(t *testing.T) { writeChart(t, "filesystem.png", RenderFilesystem) }
func TestWriteDiskIOChart(t *testing.T)     { writeChart(t, "diskio.png", RenderDiskIO) }
func TestWriteNetworkChart(t *testing.T)     { writeChart(t, "network.png", RenderNetwork) }

func writeChart(t *testing.T, filename string, render func(io.Writer) error) {
	t.Helper()
	outDir := filepath.Join("output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(outDir, filename)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := render(f); err != nil {
		t.Fatal(err)
	}
	t.Logf("written: %s", path)
}
