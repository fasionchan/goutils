/*
 * Author: fasion
 * Created time: 2024-01-05 08:53:44
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 09:05:58
 */

package netutils

import (
	"fmt"
	"testing"
)

func TestNewCookiesFromNameValueMapping(t *testing.T) {
	cookies := NewCookiesFromNameValueMapping(map[string]string{
		"name": "value",
		"foo":  "bar",
	})
	fmt.Println(cookies.String())
}
