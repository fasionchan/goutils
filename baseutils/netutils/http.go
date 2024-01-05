/*
 * Author: fasion
 * Created time: 2024-01-05 08:49:25
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 13:42:19
 */

package netutils

import (
	"net/http"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

type CookiePtr = *http.Cookie
type Cookies []*http.Cookie

func (cookies Cookies) Native() []*http.Cookie {
	return cookies
}

func (cookies Cookies) Strings() types.Strings {
	return stl.Map(cookies, CookiePtr.String)
}

func (cookies Cookies) String() string {
	return cookies.Strings().Join("; ")
}

func NewCookiesFromMap(mapping map[string]string) Cookies {
	return stl.MapMapToSlice[Cookies](mapping, func(name, value string, _ map[string]string) *http.Cookie {
		return &http.Cookie{
			Name:  name,
			Value: value,
		}
	})
}

func NewHeaderFromMap(m map[string]string) http.Header {
	return stl.MultivalueMap[http.Header](m)
}
