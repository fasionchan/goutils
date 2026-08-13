package _net

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpNets_Contains(t *testing.T) {
	assert.True(t, Ipv4Intranets.Contains(net.ParseIP("10.0.0.1")))
	assert.True(t, Ipv4Intranets.Contains(net.ParseIP("172.16.0.1")))
	assert.True(t, Ipv4Intranets.Contains(net.ParseIP("192.168.0.1")))

	assert.False(t, Ipv4Intranets.Contains(net.ParseIP("11.0.0.0")))
	assert.False(t, Ipv4Intranets.Contains(net.ParseIP("173.16.0.0")))
	assert.False(t, Ipv4Intranets.Contains(net.ParseIP("193.168.0.0")))
}

func TestGetLocalIpv4Addr(t *testing.T) {
	ip, err := GetLocalIpv4Addr()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
}