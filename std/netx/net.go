package netx

import (
	"math/rand/v2"
	"net"
)

func RandomLocalTcpAddr() (*net.TCPAddr) {
	return TcpAddrFromUint64(rand.Uint64())
}

func TcpAddrFromUint64(i uint64) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   net.IPv4(127, byte(i & 0xff), byte((i >> 8) & 0xff), byte((i >> 16) & 0xfe | 0x2)),
		Port: int((i >> 24) & 0xffff) | 0x8000,
	}
}