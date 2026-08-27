//go:build !linux

package pprofutil

import "runtime"

func readProcessMemory() (ProcessMemory, error) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ProcessMemory{Current: ms.Sys}, nil
}
