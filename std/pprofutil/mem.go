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

type MemoryThreshold struct {
	Bytes   uint64
	Percent float64
}

func (t MemoryThreshold) Validate() error {
	hasBytes := t.Bytes > 0
	hasPercent := t.Percent > 0
	switch {
	case hasBytes && hasPercent:
		return errThresholdBoth
	case !hasBytes && !hasPercent:
		return errThresholdEmpty
	case hasPercent && t.Percent > 100:
		return errThresholdPercent
	default:
		return nil
	}
}

func (t MemoryThreshold) Resolve(limit uint64) (uint64, error) {
	if err := t.Validate(); err != nil {
		return 0, err
	}
	if t.Bytes > 0 {
		return t.Bytes, nil
	}
	if isUnlimitedMemory(limit) {
		return 0, errPercentNeedsLimit
	}
	return uint64(float64(limit) * t.Percent / 100), nil
}

type HeapProfiler struct {
	path      string
	threshold MemoryThreshold
}

func (p *HeapProfiler) Profile() error {
	mem, err := ReadProcessMemory()
	if err != nil {
		return fmt.Errorf("read process memory: %w", err)
	}

	if exceeds, err := mem.Exceeds(p.threshold); err != nil {
		return fmt.Errorf("check memory exceeds: %w", err)
	} else if exceeds {
		return fmt.Errorf("memory exceeds threshold: %d > %d", mem.Current, p.threshold.Bytes)
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
		threshold: MemoryThreshold{
			Bytes:   0,
			Percent: 0,
		},
	}).Profile()
}

func WithHeapProfilerPath(path string) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.path = path
	}
}

func WithHeapProfilerThresholdBytes(bytes uint64) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.threshold.Bytes = bytes
	}
}

func WithHeapProfilerThresholdPercent(percent float64) HeapProfilerOption {
	return func(p *HeapProfiler) {
		p.threshold.Percent = percent
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

func (m ProcessMemory) Exceeds(threshold MemoryThreshold) (bool, error) {
	bound, err := threshold.Resolve(m.Limit)
	if err != nil {
		return false, err
	}
	return m.Current >= bound, nil
}
