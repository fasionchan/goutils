package netx

import (
	"math/rand/v2"
	"net"
)

func RandomLocalTcpAddr() (*net.TCPAddr) {
	r := rand.Int64()
	return &net.TCPAddr{
		IP:   net.IPv4(127, byte(r & 0xff), byte((r >> 8) & 0xff), byte((r >> 16) & 0xfe)),
		Port: int((r >> 24) & 0xffff) | 0x8000,
	}
}