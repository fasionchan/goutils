/*
 * Author: fasion
 * Created time: 2024-01-05 09:28:37
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 10:22:45
 */

package netutils

import (
	"fmt"
	"testing"
)

func TestParseUrl(t *testing.T) {
	fmt.Println(ParseAndJsonifyUrl("abc"))
	fmt.Println(ParseAndJsonifyUrl("./abc"))
	fmt.Println(ParseAndJsonifyUrl("../abc"))

	fmt.Println(ParseAndJsonifyUrl("fasionchan.com/abc"))
	fmt.Println(ParseAndJsonifyUrl("fasion@fasionchan.com/abc"))
	fmt.Println(ParseAndJsonifyUrl("fasion:pass@fasionchan.com/abc"))

	fmt.Println(ParseAndJsonifyUrl("cdn.fasionchan.com/abc"))
	fmt.Println(ParseAndJsonifyUrl("fasion@cdn.fasionchan.com/abc"))
	fmt.Println(ParseAndJsonifyUrl("fasion:pass@cdn.fasionchan.com/abc"))
}

func TestJoinRawUrl(t *testing.T) {
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", ""))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "./"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "../"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "./abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "../abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "http://cdn.fasionchan.com/abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "http://cdn.fasionchan.com/abc/../cba"))
	fmt.Println(JoinRawUrl("http://fasion@fasionchan.com/base/path", ""))
	fmt.Println(JoinRawUrl("http://fasion:pass@fasionchan.com/base/path", ""))

	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "cdn.fasionchan.com/abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "fasion@cdn.fasionchan.com/abc"))
	fmt.Println(JoinRawUrl("http://fasionchan.com/base/path", "fasion:pass@cdn.fasionchan.com/abc"))
}
