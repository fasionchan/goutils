/*
 * Author: fasion
 * Created time: 2023-11-03 11:01:21
 * Last Modified by: fasion
 * Last Modified time: 2023-11-03 11:17:07
 */

package baseutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonRawMessageMap(t *testing.T) {
	var m JsonRawMessageMap

	jsonData, err := json.Marshal(JsonMap{
		"a": "aaaa",
		"b": 2222,
		"c": JsonMap{
			"value": 3333,
		},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := json.Unmarshal(jsonData, &m); err != nil {
		t.Fatal(err)
		return
	}

	a, ok := m.Lookup("a")
	assert.True(t, ok)

	aValue, err := a.UnmarshalAsString()
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, aValue, "aaaa")

	b, ok := m.Lookup("b")
	assert.True(t, ok)

	bValue, err := b.UnmarshalAsInt()
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, bValue, 2222)

	c, ok := m.Lookup("c")
	assert.True(t, ok)

	var cValue struct {
		Value int `json:"value"`
	}
	if err := c.UnmarshalTo(&cValue); err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, cValue.Value, 3333)
}
