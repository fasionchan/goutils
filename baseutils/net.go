/*
 * Author: fasion
 * Created time: 2023-05-14 11:41:33
 * Last Modified by: fasion
 * Last Modified time: 2023-05-14 12:47:47
 */

package baseutils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseNetloc(loc string, defaultAddr string, defaultPort int) (addr string, port int, err error) {
	loc = strings.TrimSpace(loc)
	parts := strings.Split(loc, ":")

	if len(parts) > 2 {
		err = fmt.Errorf("bad smtp server: %s", loc)
		return
	}

	if len(parts) == 1 {
		addr = parts[0]
		port = defaultPort
		return
	}

	if portstr := parts[1]; portstr == "" {
		port = defaultPort
	} else if port, err = strconv.Atoi(portstr); err != nil {
		return
	}

	if addrstr := parts[0]; addrstr == "" {
		addr = defaultAddr
	} else {
		addr = addrstr
	}

	return
}
