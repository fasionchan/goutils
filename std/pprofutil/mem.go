package pprofutil

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/fasionchan/goutils/stl"
)

var (
	errThresholdEmpty    = errors.New("memory threshold: Bytes or Percent is required")
	errThresholdBoth     = errors.New("memory threshold: set either Bytes or Percent, not both")
	errThresholdPercent  = errors.New("memory threshold: Percent must be in (0, 100]")
	errPercentNeedsLimit = errors.New("memory threshold: percent requires a finite memory limit (k8s/cgroup)")
)

// MemoryThreshold 内存阈值：绝对值（字节）或相对内存上限的百分比。两者都设置时，有一个满足即可。
// 百分比相对 ReadProcessMemory 得到的 Limit（优先 k8s/cgroup memory limit）。
type MemoryThreshold struct {
	Bytes   *uint64
	Percent *float64
}

func (t MemoryThreshold) IsEnabled() bool {
	return t.Bytes != nil || t.Percent != nil
}

func (t MemoryThreshold) Exceeds(current uint64, limit uint64) bool {
	return t.BytesExceeds(current) || t.PercentExceeds(current, limit)
}

func (t MemoryThreshold) BytesExceeds(current uint64) bool {
	bytes := t.Bytes
	if bytes == nil {
		return false
	}

	return current >= *bytes
}

func (t MemoryThreshold) PercentExceeds(current uint64, limit uint64) bool {
	percent := t.Percent
	if percent == nil {
		return false
	}

	return current >= uint64(float64(limit) * *percent / 100)
}

type HeapProfiler struct {
	path      string
	threshold MemoryThreshold
}

func (p *HeapProfiler) Profile() error {
	if p.threshold.IsEnabled() {
		mem, err := ReadProcessMemory()
		if err != nil {
			return fmt.Errorf("read process memory: %w", err)
		}

		if !mem.Exceeds(p.threshold) {
			return nil
		}
	}

	path := p.path
	if path == "" {
		return fmt.Errorf("heap profile path is required")
	}

	f, err := os.Create(p.path)
	if err != nil {
		return fmt.Errorf("create heap profile file %q: %w", path, err)
	}
	defer f.Close()

	runtime.GC()
	if err := pprof.Lookup("heap").WriteTo(f, 0); err != nil {
		return fmt.Errorf("write heap profile: %w", err)
	}

	return nil
}

type HeapProfilerOption = stl.Option[*HeapProfiler]

func ProfileHeap(opts ...HeapProfilerOption) error {
	return stl.NewOptions(opts...).Apply(&HeapProfiler{
		path: "heap.pprof",
	}).Profile()
}

func WithHeapProfilerPath(path string) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.path = path
	}
}

func WithHeapProfilerThresholdBytes(bytes uint64) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.threshold.Bytes = &bytes
	}
}

func WithHeapProfilerThresholdPercent(percent float64) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.threshold.Percent = &percent
	}
}

// ProcessMemory 当前进程/容器内存用量与可用于百分比的上限。
type ProcessMemory struct {
	Current uint64
	Limit   uint64
}

func ReadProcessMemory() (ProcessMemory, error) {
	return readProcessMemory()
}

func (m ProcessMemory) Exceeds(threshold MemoryThreshold) bool {
	return threshold.Exceeds(m.Current, m.Limit)
}
