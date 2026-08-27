package pprofutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	unlimitedMemoryLimit = uint64(1) << 62
)

func isUnlimitedMemory(limit uint64) bool {
	return limit == 0 || limit >= unlimitedMemoryLimit
}

func parseCgroupLimit(s string) (uint64, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "max" {
		return 0, false
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil || isUnlimitedMemory(n) {
		return 0, false
	}
	return n, true
}

func parseUintText(s string) (uint64, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse uint %q: %w", strings.TrimSpace(s), err)
	}
	return n, nil
}

func parseVmRSS(status string) (uint64, error) {
	for _, line := range strings.Split(status, "\n") {
		after, ok := strings.CutPrefix(line, "VmRSS:")
		if !ok {
			continue
		}
		fields := strings.Fields(after)
		if len(fields) == 0 {
			return 0, fmt.Errorf("parse VmRSS: empty value")
		}
		kb, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse VmRSS: %w", err)
		}
		return kb * 1024, nil
	}
	return 0, errors.New("VmRSS not found")
}
