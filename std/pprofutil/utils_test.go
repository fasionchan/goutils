package pprofutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVmRSS(t *testing.T) {
	status := "Name:\tapp\nVmRSS:\t   2048 kB\nVmSize:\t   4096 kB\n"
	got, err := parseVmRSS(status)
	require.NoError(t, err)
	assert.Equal(t, uint64(2048*1024), got)

	_, err = parseVmRSS("Name:\tapp\n")
	assert.Error(t, err)
}

func TestParseCgroupLimit(t *testing.T) {
	limit, ok := parseCgroupLimit("1073741824\n")
	assert.True(t, ok)
	assert.Equal(t, uint64(1073741824), limit)

	_, ok = parseCgroupLimit("max\n")
	assert.False(t, ok)

	_, ok = parseCgroupLimit("9223372036854771712")
	assert.False(t, ok)
}
