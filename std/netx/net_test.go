package netx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpAddrFromInt64(t *testing.T) {
	assert.Equal(t, TcpAddrFromUint64(1 | (2<<8) | (3<<16) | (4<<24)).String(), "127.1.2.2:32772")
	assert.Equal(t, TcpAddrFromUint64(1 | (2<<8) | (5<<16) | (4<<24)).String(), "127.1.2.6:32772")
	assert.Equal(t, TcpAddrFromUint64(1| (2<<8) | (0<<16) | (3<<24)).String(), "127.1.2.2:32771")
	assert.Equal(t, TcpAddrFromUint64(1| (2<<8) | (255<<16) | (3<<24)).String(), "127.1.2.254:32771")
}


func TestRandomLocalTcpAddr(t *testing.T) {
	addr := RandomLocalTcpAddr()
	fmt.Println(addr)

	addr2 := RandomLocalTcpAddr()
	fmt.Println(addr2)

	assert.NotEqual(t, addr.IP, addr2.IP)
	assert.NotEqual(t, addr.Port, addr2.Port)
}