/*
 * Author: fasion
 * Created time: 2025-12-07 13:22:36
 * Last Modified by: fasion
 * Last Modified time: 2025-12-08 13:24:21
 */

package entity

import (
	"fmt"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFlattenDataAttributes(t *testing.T) {
	attrs := FlattenDataAttributes(map[string]any{
		"a": 1,
		"b": "2",
		"c": []int{1, 2, 3},
		"d": struct {
			E        int
			F        string
			G        []int
			OmitZero int `attr:"OmitZero,omitzero"`
			H        int `attr:"-"`
		}{
			E:        4,
			F:        "5",
			G:        []int{6, 7, 8},
			OmitZero: 0,
			H:        9,
		},
		"o": primitive.NewObjectID(),
		"x": []struct {
			Y int
			Z int
		}{
			{
				Y: 10,
				Z: 11,
			},
			{
				Y: 12,
				Z: 13,
			},
		},
	})

	attrs.SortByName().Print()
	fmt.Println()

	for _, attr := range attrs {
		fmt.Println(attr.Name, IndexPattern.ReplaceAllString(attr.Name, ""))
	}
}

func init() {
	RegisterAtomicAttrType(primitive.NilObjectID)
	RegisterData2Attr(func(data primitive.ObjectID) *Attribute {
		return &Attribute{Value: strings.ToUpper(data.Hex())}
	})
}
