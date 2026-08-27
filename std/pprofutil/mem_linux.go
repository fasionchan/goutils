//go:build linux

package pprofutil

import (
	"fmt"
	"os"
	"runtime"
)

const (
	cgroupV2Current = "/sys/fs/cgroup/memory.current"
	cgroupV2Max     = "/sys/fs/cgroup/memory.max"
	cgroupV1Usage   = "/sys/fs/cgroup/memory/memory.usage_in_bytes"
	cgroupV1Limit   = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	procSelfStatus  = "/proc/self/status"
)

func readProcessMemory() (ProcessMemory, error) {
	if mem, ok := readCgroupMemory(); ok && mem.Limit > 0 {
		return mem, nil
	}

	rss, err := readProcVmRSS()
	if err != nil {
		return fallbackGoMemory()
	}
	return ProcessMemory{Current: rss}, nil
}

func readCgroupMemory() (ProcessMemory, bool) {
	if mem, ok := readCgroupPair(cgroupV2Current, cgroupV2Max); ok {
		return mem, true
	}
	return readCgroupPair(cgroupV1Usage, cgroupV1Limit)
}

func readCgroupPair(currentPath, limitPath string) (ProcessMemory, bool) {
	currentRaw, err := os.ReadFile(currentPath)
	if err != nil {
		return ProcessMemory{}, false
	}
	current, err := parseUintText(string(currentRaw))
	if err != nil {
		return ProcessMemory{}, false
	}

	limitRaw, err := os.ReadFile(limitPath)
	if err != nil {
		return ProcessMemory{Current: current}, true
	}
	limit, finite := parseCgroupLimit(string(limitRaw))
	if !finite {
		return ProcessMemory{Current: current}, true
	}
	return ProcessMemory{Current: current, Limit: limit}, true
}

func readProcVmRSS() (uint64, error) {
	data, err := os.ReadFile(procSelfStatus)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", procSelfStatus, err)
	}
	return parseVmRSS(string(data))
}

func fallbackGoMemory() (ProcessMemory, error) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ProcessMemory{Current: ms.Sys}, nil
}
