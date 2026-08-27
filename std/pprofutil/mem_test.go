package pprofutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MemoryThresholdValidateTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	wantErr   error
}

func (tc MemoryThresholdValidateTestCase) Run(t *testing.T) {
	err := tc.threshold.Validate()
	if tc.wantErr == nil {
		assert.NoError(t, err)
		return
	}
	assert.ErrorIs(t, err, tc.wantErr)
}

func TestMemoryThresholdValidate(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdValidateTestCase{
		{
			TestCaseName: "bytes",
			threshold:    MemoryThreshold{Bytes: 1024},
		},
		{
			TestCaseName: "percent",
			threshold:    MemoryThreshold{Percent: 80},
		},
		{
			TestCaseName: "percent 100",
			threshold:    MemoryThreshold{Percent: 100},
		},
		{
			TestCaseName: "empty",
			threshold:    MemoryThreshold{},
			wantErr:      errThresholdEmpty,
		},
		{
			TestCaseName: "both",
			threshold:    MemoryThreshold{Bytes: 1, Percent: 1},
			wantErr:      errThresholdBoth,
		},
		{
			TestCaseName: "percent too large",
			threshold:    MemoryThreshold{Percent: 101},
			wantErr:      errThresholdPercent,
		},
	})
}

type MemoryThresholdResolveTestCase struct {
	_testing.TestCaseName
	threshold MemoryThreshold
	limit     uint64
	want      uint64
	wantErr   error
}

func (tc MemoryThresholdResolveTestCase) Run(t *testing.T) {
	got, err := tc.threshold.Resolve(tc.limit)
	if tc.wantErr != nil {
		assert.ErrorIs(t, err, tc.wantErr)
		assert.Equal(t, uint64(0), got)
		return
	}
	require.NoError(t, err)
	assert.Equal(t, tc.want, got)
}

func TestMemoryThresholdResolve(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []MemoryThresholdResolveTestCase{
		{
			TestCaseName: "bytes ignores limit",
			threshold:    MemoryThreshold{Bytes: 4096},
			want:         4096,
		},
		{
			TestCaseName: "percent of limit",
			threshold:    MemoryThreshold{Percent: 50},
			limit:        1000,
			want:         500,
		},
		{
			TestCaseName: "percent of zero limit",
			threshold:    MemoryThreshold{Percent: 80},
			wantErr:      errPercentNeedsLimit,
		},
		{
			TestCaseName: "percent of unlimited",
			threshold:    MemoryThreshold{Percent: 80},
			limit:        unlimitedMemoryLimit,
			wantErr:      errPercentNeedsLimit,
		},
		{
			TestCaseName: "invalid empty",
			threshold:    MemoryThreshold{},
			wantErr:      errThresholdEmpty,
		},
	})
}

type ProcessMemoryExceedsTestCase struct {
	_testing.TestCaseName
	mem       ProcessMemory
	threshold MemoryThreshold
	want      bool
	wantErr   error
}

func (tc ProcessMemoryExceedsTestCase) Run(t *testing.T) {
	got, err := tc.mem.Exceeds(tc.threshold)
	if tc.wantErr != nil {
		assert.ErrorIs(t, err, tc.wantErr)
		assert.False(t, got)
		return
	}
	require.NoError(t, err)
	assert.Equal(t, tc.want, got)
}

func TestProcessMemoryExceeds(t *testing.T) {
	mem := ProcessMemory{Current: 800, Limit: 1000}
	_testing.TypedRunNamedTestCases(t, []ProcessMemoryExceedsTestCase{
		{
			TestCaseName: "bytes equal is over",
			mem:          mem,
			threshold:    MemoryThreshold{Bytes: 800},
			want:         true,
		},
		{
			TestCaseName: "bytes under",
			mem:          mem,
			threshold:    MemoryThreshold{Bytes: 801},
			want:         false,
		},
		{
			TestCaseName: "percent equal is over",
			mem:          mem,
			threshold:    MemoryThreshold{Percent: 80},
			want:         true,
		},
		{
			TestCaseName: "percent under",
			mem:          mem,
			threshold:    MemoryThreshold{Percent: 80.1},
			want:         false,
		},
		{
			TestCaseName: "invalid threshold",
			mem:          mem,
			threshold:    MemoryThreshold{},
			wantErr:      errThresholdEmpty,
		},
		{
			TestCaseName: "percent needs limit",
			mem:          ProcessMemory{Current: 800},
			threshold:    MemoryThreshold{Percent: 80},
			wantErr:      errPercentNeedsLimit,
		},
	})
}

type ProfileHeapTestCase struct {
	_testing.TestCaseName
	emptyPath        bool
	seedExisting     bool
	thresholdBytes   uint64
	thresholdPercent float64
	wantErrIs        error
	wantExceeds      bool
	wantDumped       bool
}

func (tc ProfileHeapTestCase) Run(t *testing.T) {
	path := filepath.Join(t.TempDir(), "heap.pprof")
	opts := []HeapProfilerOption{WithHeapProfilerPath(path)}
	if tc.emptyPath {
		opts = []HeapProfilerOption{WithHeapProfilerPath("")}
		path = ""
	}
	if tc.thresholdBytes > 0 {
		opts = append(opts, WithHeapProfilerThresholdBytes(tc.thresholdBytes))
	}
	if tc.thresholdPercent > 0 {
		opts = append(opts, WithHeapProfilerThresholdPercent(tc.thresholdPercent))
	}
	if tc.seedExisting && path != "" {
		require.NoError(t, os.WriteFile(path, []byte("old"), 0o644))
	}

	err := ProfileHeap(opts...)
	switch {
	case tc.wantErrIs != nil:
		assert.ErrorIs(t, err, tc.wantErrIs)
	case tc.wantExceeds:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "memory exceeds threshold")
	case tc.wantDumped:
		require.NoError(t, err)
	default:
		require.Error(t, err)
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
			TestCaseName:   "dump under high bytes threshold",
			thresholdBytes: ^uint64(0),
			wantDumped:     true,
		},
		{
			TestCaseName:   "overwrite existing profile",
			seedExisting:   true,
			thresholdBytes: ^uint64(0),
			wantDumped:     true,
		},
		{
			TestCaseName:   "refuse when over bytes threshold",
			thresholdBytes: 1,
			wantExceeds:    true,
		},
		{
			TestCaseName:   "refuse does not overwrite",
			seedExisting:   true,
			thresholdBytes: 1,
			wantExceeds:    true,
		},
		{
			TestCaseName: "empty threshold",
			wantErrIs:    errThresholdEmpty,
		},
		{
			TestCaseName:   "empty path",
			emptyPath:      true,
			thresholdBytes: ^uint64(0),
		},
		{
			TestCaseName:     "both thresholds",
			thresholdBytes:   1024,
			thresholdPercent: 80,
			wantErrIs:        errThresholdBoth,
		},
	})
}

func TestHeapProfilerProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "heap.pprof")
	p := &HeapProfiler{
		path:      path,
		threshold: MemoryThreshold{Bytes: ^uint64(0)},
	}
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
