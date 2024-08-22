/*
 * Author: fasion
 * Created time: 2024-08-22 14:19:40
 * Last Modified by: fasion
 * Last Modified time: 2024-08-22 14:22:50
 */

package queryutils

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/types"
)

func TestNewSubdataSetinExpression(t *testing.T) {
	fmt.Println(NewSubdataSetinExpression(types.NewStrings("a", "b", "c"), "D"))
}

func TestNewSubdataSetinExpressionX(t *testing.T) {
	fmt.Println(NewSubdataSetinExpressionX("D", "a", "b", "c"))
}
