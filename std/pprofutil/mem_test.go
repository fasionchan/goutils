package pprofutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/fasionchan/goutils/stl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MemoryThresholdIsEnabledTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	want      bool
}

func (tc MemoryThresholdIsEnabledTestCase) Run(t *testing.T) {
	assert.Equal(t, tc.want, tc.threshold.IsEnabled())
}

func TestMemoryThresholdIsEnabled(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdIsEnabledTestCase{
		{
			TestCaseName: "bytes",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(1024))},
			want:         true,
		},
		{
			TestCaseName: "percent",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.0)},
			want:         true,
		},
		{
			TestCaseName: "both",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(1)), Percent: stl.AddrOf(1.0)},
			want:         true,
		},
		{
			TestCaseName: "empty",
			threshold:    MemoryThreshold{},
		},
	})
}

type MemoryThresholdExceedsTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	current   uint64
	limit     uint64
	want      bool
}

func (tc MemoryThresholdExceedsTestCase) Run(t *testing.T) {
	assert.Equal(t, tc.want, tc.threshold.Exceeds(tc.current, tc.limit))
	assert.Equal(t, tc.want, ProcessMemory{Current: tc.current, Limit: tc.limit}.Exceeds(tc.threshold))
}

func TestMemoryThresholdExceeds(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdExceedsTestCase{
		{
			TestCaseName: "bytes equal is over",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(800))},
			current:      800,
			want:         true,
		},
		{
			TestCaseName: "bytes under",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(801))},
			current:      800,
		},
		{
			TestCaseName: "bytes ignores limit",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(800))},
			current:      800,
			limit:        1,
			want:         true,
		},
		{
			TestCaseName: "percent equal is over",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.0)},
			current:      800,
			limit:        1000,
			want:         true,
		},
		{
			TestCaseName: "percent under",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.1)},
			current:      800,
			limit:        1000,
		},
		{
			TestCaseName: "percent of zero limit always over",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.0)},
			current:      800,
			want:         true,
		},
		{
			TestCaseName: "empty never exceeds",
			threshold:    MemoryThreshold{},
			current:      800,
			limit:        1000,
		},
		{
			TestCaseName: "both bytes hits",
			threshold: MemoryThreshold{
				Bytes:   stl.AddrOf(uint64(800)),
				Percent: stl.AddrOf(90.0),
			},
			current: 800,
			limit:   1000,
			want:    true,
		},
		{
			TestCaseName: "both percent hits",
			threshold: MemoryThreshold{
				Bytes:   stl.AddrOf(uint64(900)),
				Percent: stl.AddrOf(80.0),
			},
			current: 800,
			limit:   1000,
			want:    true,
		},
		{
			TestCaseName: "both miss",
			threshold: MemoryThreshold{
				Bytes:   stl.AddrOf(uint64(801)),
				Percent: stl.AddrOf(80.1),
			},
			current: 800,
			limit:   1000,
		},
	})
}

type MemoryThresholdBytesExceedsTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	current   uint64
	want      bool
}

func (tc MemoryThresholdBytesExceedsTestCase) Run(t *testing.T) {
	assert.Equal(t, tc.want, tc.threshold.BytesExceeds(tc.current))
}

func TestMemoryThresholdBytesExceeds(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdBytesExceedsTestCase{
		{
			TestCaseName: "nil",
			current:      800,
		},
		{
			TestCaseName: "equal",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(800))},
			current:      800,
			want:         true,
		},
		{
			TestCaseName: "under",
			threshold:    MemoryThreshold{Bytes: stl.AddrOf(uint64(801))},
			current:      800,
		},
	})
}

type MemoryThresholdPercentExceedsTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	current   uint64
	limit     uint64
	want      bool
}

func (tc MemoryThresholdPercentExceedsTestCase) Run(t *testing.T) {
	assert.Equal(t, tc.want, tc.threshold.PercentExceeds(tc.current, tc.limit))
}

func TestMemoryThresholdPercentExceeds(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdPercentExceedsTestCase{
		{
			TestCaseName: "nil",
			current:      800,
			limit:        1000,
		},
		{
			TestCaseName: "equal",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.0)},
			current:      800,
			limit:        1000,
			want:         true,
		},
		{
			TestCaseName: "under",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.1)},
			current:      800,
			limit:        1000,
		},
		{
			TestCaseName: "zero limit",
			threshold:    MemoryThreshold{Percent: stl.AddrOf(80.0)},
			current:      1,
			want:         true,
		},
	})
}

type ProfileHeapTestCase struct {
	_testing.TestCaseName
	emptyPath      bool
	seedExisting   bool
	bytes          *uint64
	wantErrContain string
	wantDumped     bool
}

func (tc ProfileHeapTestCase) Run(t *testing.T) {
	path := filepath.Join(t.TempDir(), "heap.pprof")
	opts := []HeapProfilerOption{WithHeapProfilerPath(path)}
	if tc.emptyPath {
		opts = []HeapProfilerOption{WithHeapProfilerPath("")}
		path = ""
	}
	if tc.bytes != nil {
		opts = append(opts, WithHeapProfilerThresholdBytes(*tc.bytes))
	}
	if tc.seedExisting && path != "" {
		require.NoError(t, os.WriteFile(path, []byte("old"), 0o644))
	}

	err := ProfileHeap(opts...)
	if tc.wantErrContain != "" {
		require.Error(t, err)
		assert.Contains(t, err.Error(), tc.wantErrContain)
	} else {
		require.NoError(t, err)
	}

	if path == "" {
		return
	}
	data, readErr := os.ReadFile(path)
	if tc.wantDumped {
		require.NoError(t, readErr)
		require.Greater(t, len(data), 2)
		assert.Equal(t, []byte{0x1f, 0x8b}, data[:2])
		return
	}
	if tc.seedExisting {
		assert.Equal(t, []byte("old"), data)
		return
	}
	assert.ErrorIs(t, readErr, os.ErrNotExist)
}

func TestProfileHeap(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []ProfileHeapTestCase{
		{
			TestCaseName: "dump without threshold",
			wantDumped:   true,
		},
		{
			TestCaseName: "overwrite existing profile",
			seedExisting: true,
			wantDumped:   true,
		},
		{
			TestCaseName: "dump when over bytes threshold",
			bytes:        stl.AddrOf(uint64(1)),
			wantDumped:   true,
		},
		{
			TestCaseName: "skip when under bytes threshold",
			bytes:        stl.AddrOf(^uint64(0)),
		},
		{
			TestCaseName: "skip does not overwrite",
			seedExisting: true,
			bytes:        stl.AddrOf(^uint64(0)),
		},
		{
			TestCaseName:   "empty path",
			emptyPath:      true,
			wantErrContain: "heap profile path is required",
		},
		{
			TestCaseName: "empty path skipped under threshold",
			emptyPath:    true,
			bytes:        stl.AddrOf(^uint64(0)),
		},
	})
}

func TestHeapProfilerProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "heap.pprof")
	p := &HeapProfiler{path: path}
	require.NoError(t, p.Profile())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Greater(t, len(data), 2)
	assert.Equal(t, []byte{0x1f, 0x8b}, data[:2])
}

func TestReadProcessMemory(t *testing.T) {
	mem, err := ReadProcessMemory()
	require.NoError(t, err)
	assert.Greater(t, mem.Current, uint64(0))
}
