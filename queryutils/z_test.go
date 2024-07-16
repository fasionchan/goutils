/*
 * Author: fasion
 * Created time: 2023-05-24 16:14:26
 * Last Modified by: fasion
 * Last Modified time: 2024-07-16 09:57:00
 */

package queryutils

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/types"
)

func TestCompiling(t *testing.T) {
	a := "abc"
	fmt.Printf("%v\n", a[len(a):])

}

func TestParseSetinExpression(t *testing.T) {
	fmt.Println(ParseSetinExpression("A"))
	fmt.Println(ParseSetinExpression("(A)"))
	fmt.Println(ParseSetinExpression("((A))"))
	fmt.Println(ParseSetinExpression("(A)(B)"))
	fmt.Println(ParseSetinExpression("((A)(B))"))
	fmt.Println(ParseSetinExpression(""))
	fmt.Println(ParseSetinExpression("   "))
	fmt.Println(ParseSetinExpression("(A)B(C)D(E)F"))
	fmt.Println(ParseSetinExpression("(A)(B)C"))
	fmt.Println(ParseSetinExpression("(A)("))
	fmt.Println(ParseSetinExpression("(A)(C"))
	fmt.Println(ParseSetinExpression("(A))"))
	fmt.Println(ParseSetinExpression("(A))("))
	fmt.Println(ParseSetinExpression("("))
	fmt.Println(ParseSetinExpression("()"))
	fmt.Println(ParseSetinExpression("(()"))
	fmt.Println(ParseSetinExpression(")()"))
	fmt.Println(ParseSetinExpression(")())"))
}

func TestEssentialDataTypeIdent(t *testing.T) {
	fmt.Println(EssentialDataTypeIdent(types.Strings{}))
	fmt.Println(EssentialDataTypeIdent(types.NewStrings().Empty))
	fmt.Println(EssentialDataTypeIdent(baseutils.BadTypeError{}))
	fmt.Println(EssentialDataTypeIdent(&baseutils.BadTypeError{}))
	fmt.Println(EssentialDataTypeIdent(&Setiner[[]*int, []int, *int, int]{}))
}
