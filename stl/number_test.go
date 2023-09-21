/*
 * Author: fasion
 * Created time: 2023-09-21 11:27:45
 * Last Modified by: fasion
 * Last Modified time: 2023-09-21 11:36:34
 */

package stl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	assert.Equal(t, Max([]int{}, 1), 1)
	assert.Equal(t, Max([]int{2}, 1), 2)
	assert.Equal(t, Max([]int{2, 3, 4, 5}, 1), 5)

	assert.Equal(t, Min([]int{}, 1), 1)
	assert.Equal(t, Min([]int{2}, 1), 2)
	assert.Equal(t, Min([]int{2, 3, 4, 5}, 1), 2)
	assert.Equal(t, Min([]int{2, 3, 4, 5, 0}, 1), 0)

	assert.Equal(t, Sum([]int{}, 1), 1)
	assert.Equal(t, Sum([]int{2}, 1), 3)
	assert.Equal(t, Sum([]int{2, 3, 4, 5}, 1), 15)
}
