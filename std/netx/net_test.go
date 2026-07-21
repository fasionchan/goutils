package netx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRandomLocalTcpAddr(t *testing.T) {
	addr := RandomLocalTcpAddr()
	fmt.Println(addr)

	addr2 := RandomLocalTcpAddr()
	fmt.Println(addr2)

	assert.NotEqual(t, addr.IP, addr2.IP)
	assert.NotEqual(t, addr.Port, addr2.Port)
}