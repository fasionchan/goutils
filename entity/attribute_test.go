/*
 * Author: fasion
 * Created time: 2025-12-07 13:22:36
 * Last Modified by: fasion
 * Last Modified time: 2025-12-08 00:52:53
 */

package entity

import (
	"fmt"
	"testing"
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
	for _, attr := range attrs {
		fmt.Println(attr.Name, IndexPattern.ReplaceAllString(attr.Name, ""))
	}
}
