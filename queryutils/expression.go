/*
 * Author: fasion
 * Created time: 2024-06-08 01:03:36
 * Last Modified by: fasion
 * Last Modified time: 2024-06-08 08:24:52
 */

package queryutils

import (
	"strings"

	"github.com/fasionchan/goutils/baseutils"
)

func parseSetinExpression(setin string) (result []string, err error) {
Start:
	n := len(setin)
	if n == 0 {
		return
	}

	brackets := 0
	bracketStart := -1

	for i := 0; i < n; i++ {
		switch setin[i] {
		case '(':
			if brackets == 0 {
				if i > 0 {
					result = append(result, setin[:i])
				}
				bracketStart = i
			}
			brackets++
		case ')':
			brackets--
			if brackets == 0 {
				result = append(result, setin[bracketStart+1:i])
				setin = setin[i+1:]
				goto Start
			} else if brackets < 0 {
				return nil, baseutils.NewBadTypeError("Setin", setin)
			}
		}
	}

	if brackets != 0 {
		return nil, baseutils.NewBadTypeError("Setin", setin)
	}

	result = append(result, setin)

	return
}

func ParseSetinExpression(setin string) ([]string, error) {
	return parseSetinExpression(strings.TrimSpace(setin))
}
