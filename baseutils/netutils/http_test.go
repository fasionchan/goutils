/*
 * Author: fasion
 * Created time: 2024-01-05 08:53:44
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 13:24:05
 */

package netutils

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewCookiesFromNameValueMapping(t *testing.T) {
	cookies := NewCookiesFromMap(map[string]string{
		"name": "value",
		"foo":  "bar",
	})
	fmt.Println(cookies.String())
}

func TestNewHeaderFromMap(t *testing.T) {
	headers := NewHeaderFromMap(map[string]string{
		"name": "value",
		"foo":  "bar",
	})
	fmt.Println(headers)
}

func TestXxx(t *testing.T) {
	http.Header(nil).Set("", "")
}
