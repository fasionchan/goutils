package runtimeutil

import "runtime"

type MemStats runtime.MemStats

func (stats *MemStats) Native() *runtime.MemStats {
	return (*runtime.MemStats)(stats)
}

func (stats *MemStats) Sub(other *MemStats) *MemStats {
	result := *stats
	result.Alloc -= other.Alloc
	result.TotalAlloc -= other.TotalAlloc
	result.Sys -= other.Sys
	result.Lookups -= other.Lookups
	result.Mallocs -= other.Mallocs
	result.Frees -= other.Frees
	result.HeapAlloc -= other.HeapAlloc
	result.HeapSys -= other.HeapSys
	result.HeapIdle -= other.HeapIdle
	result.HeapInuse -= other.HeapInuse
	result.HeapReleased -= other.HeapReleased
	result.HeapObjects -= other.HeapObjects
	result.StackInuse -= other.StackInuse
	result.StackSys -= other.StackSys
	result.MSpanInuse -= other.MSpanInuse
	result.MSpanSys -= other.MSpanSys
	result.MCacheInuse -= other.MCacheInuse
	result.MCacheSys -= other.MCacheSys
	result.BuckHashSys -= other.BuckHashSys
	result.GCSys -= other.GCSys
	result.OtherSys -= other.OtherSys
	result.NextGC -= other.NextGC
	result.LastGC -= other.LastGC
	result.PauseTotalNs -= other.PauseTotalNs
	result.NumGC -= other.NumGC
	result.NumForcedGC -= other.NumForcedGC
	result.GCCPUFraction -= other.GCCPUFraction

	for i := range result.BySize {
		result.BySize[i].Mallocs -= other.BySize[i].Mallocs
		result.BySize[i].Frees -= other.BySize[i].Frees
	}

	return &result
}

func (stats *MemStats) Delta(fn func(delta *MemStats)) {
	fn(GetMemStats().Sub(stats))
}

func GetMemStats() *MemStats {
	var result runtime.MemStats
	runtime.ReadMemStats(&result)
	return (*MemStats)(&result)
}
