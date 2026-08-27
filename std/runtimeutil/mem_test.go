package runtimeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStats(t *testing.T) {
	defer GetMemStats().Delta(func(delta *MemStats) {
		t.Logf("delta: %+v", delta)
	})
}

func TestMemStatsSub(t *testing.T) {
	later := &MemStats{
		Alloc:         110,
		TotalAlloc:    210,
		Sys:           310,
		Lookups:       4,
		Mallocs:       50,
		Frees:         20,
		HeapAlloc:     111,
		HeapSys:       211,
		HeapIdle:      12,
		HeapInuse:     13,
		HeapReleased:  14,
		HeapObjects:   15,
		StackInuse:    21,
		StackSys:      22,
		MSpanInuse:    31,
		MSpanSys:      32,
		MCacheInuse:   41,
		MCacheSys:     42,
		BuckHashSys:   51,
		GCSys:         61,
		OtherSys:      71,
		NextGC:        81,
		LastGC:        91,
		PauseTotalNs:  101,
		NumGC:         8,
		NumForcedGC:   3,
		GCCPUFraction: 0.5,
		EnableGC:      true,
		DebugGC:       true,
	}
	later.PauseNs[0] = 1000
	later.PauseEnd[0] = 2000
	later.BySize[0].Size = 8
	later.BySize[0].Mallocs = 9
	later.BySize[0].Frees = 6

	earlier := &MemStats{
		Alloc:         10,
		TotalAlloc:    10,
		Sys:           10,
		Lookups:       1,
		Mallocs:       5,
		Frees:         2,
		HeapAlloc:     11,
		HeapSys:       11,
		HeapIdle:      2,
		HeapInuse:     3,
		HeapReleased:  4,
		HeapObjects:   5,
		StackInuse:    1,
		StackSys:      2,
		MSpanInuse:    1,
		MSpanSys:      2,
		MCacheInuse:   1,
		MCacheSys:     2,
		BuckHashSys:   1,
		GCSys:         1,
		OtherSys:      1,
		NextGC:        1,
		LastGC:        1,
		PauseTotalNs:  1,
		NumGC:         2,
		NumForcedGC:   1,
		GCCPUFraction: 0.1,
		EnableGC:      false,
		DebugGC:       false,
	}
	earlier.PauseNs[0] = 100
	earlier.PauseEnd[0] = 200
	earlier.BySize[0].Size = 8
	earlier.BySize[0].Mallocs = 3
	earlier.BySize[0].Frees = 1

	delta := later.Sub(earlier)
	assert.Equal(t, uint64(100), delta.Alloc)
	assert.Equal(t, uint64(200), delta.TotalAlloc)
	assert.Equal(t, uint64(300), delta.Sys)
	assert.Equal(t, uint64(3), delta.Lookups)
	assert.Equal(t, uint64(45), delta.Mallocs)
	assert.Equal(t, uint64(18), delta.Frees)
	assert.Equal(t, uint64(100), delta.HeapAlloc)
	assert.Equal(t, uint64(200), delta.HeapSys)
	assert.Equal(t, uint64(10), delta.HeapIdle)
	assert.Equal(t, uint64(10), delta.HeapInuse)
	assert.Equal(t, uint64(10), delta.HeapReleased)
	assert.Equal(t, uint64(10), delta.HeapObjects)
	assert.Equal(t, uint64(20), delta.StackInuse)
	assert.Equal(t, uint64(20), delta.StackSys)
	assert.Equal(t, uint64(30), delta.MSpanInuse)
	assert.Equal(t, uint64(30), delta.MSpanSys)
	assert.Equal(t, uint64(40), delta.MCacheInuse)
	assert.Equal(t, uint64(40), delta.MCacheSys)
	assert.Equal(t, uint64(50), delta.BuckHashSys)
	assert.Equal(t, uint64(60), delta.GCSys)
	assert.Equal(t, uint64(70), delta.OtherSys)
	assert.Equal(t, uint64(80), delta.NextGC)
	assert.Equal(t, uint64(90), delta.LastGC)
	assert.Equal(t, uint64(100), delta.PauseTotalNs)
	assert.Equal(t, uint32(6), delta.NumGC)
	assert.Equal(t, uint32(2), delta.NumForcedGC)
	assert.InDelta(t, 0.4, delta.GCCPUFraction, 1e-9)
	assert.Equal(t, uint64(6), delta.BySize[0].Mallocs)
	assert.Equal(t, uint64(5), delta.BySize[0].Frees)
	assert.Equal(t, uint32(8), delta.BySize[0].Size)
	assert.Equal(t, later.PauseNs, delta.PauseNs)
	assert.Equal(t, later.PauseEnd, delta.PauseEnd)
	assert.True(t, delta.EnableGC)
	assert.True(t, delta.DebugGC)
}
