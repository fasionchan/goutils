/*
 * Author: fasion
 * Created time: 2022-11-14 11:31:44
 * Last Modified by: fasion
 * Last Modified time: 2023-12-18 10:26:45
 */

package stl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testNumbers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func TestDivide(t *testing.T) {
	case1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(Divide(case1, 1, false))
	fmt.Println(Divide(case1, 2, false))
	fmt.Println(Divide(case1[1:], 2, false))
	fmt.Println(Divide(case1, 9, false))
	fmt.Println(Divide(case1, 10, false))
	fmt.Println(Divide(case1, 11, false))
	fmt.Println(Divide(case1, 0, false))
	fmt.Println(Divide(case1, -1, false))
}

func TestSlice(t *testing.T) {
	assert.Equal(t, true, AnyMatch(testNumbers, func(i int) bool {
		return i > 0
	}))
	assert.Equal(t, false, AnyMatch(testNumbers, func(i int) bool {
		return i > 10
	}))

	assert.Equal(t, true, AllMatch(testNumbers, func(i int) bool {
		return i > -1
	}))
	assert.Equal(t, false, AllMatch(testNumbers, func(i int) bool {
		return i > 0
	}))

	assert.Equal(t, []int{0, 2, 4}, Map([]int{0, 1, 2}, func(i int) int { return i * 2 }))

	assert.Equal(t, 45, Reduce(testNumbers, func(previous int, current int) int {
		return current + previous
	}, 0))

	assert.Equal(t, ConcatSlices([]int{1, 2}, []int{3, 4}), []int{1, 2, 3, 4})
}

func TestFilterDemo(t *testing.T) {
	// 应用实例：对数列整体乘以2
	var testNumbers = []int{1, 2, 3}

	var odd = func(i int) bool {
		return i%2 == 1
	}

	fmt.Println(Filter(testNumbers, odd)) // 输出：[1 3]
}

func TestFindDemo(t *testing.T) {
	// 应用实例：找出第一个奇数
	var testNumbers = []int{0, 1, 2, 3}

	var odd = func(i int) bool {
		return i%2 == 1
	}

	fmt.Println(FindFirst(testNumbers, odd)) // 输出：1 true
}

func TestMapDemo(t *testing.T) {
	// 应用实例：对数列整体乘以2
	var testNumbers = []int{1, 2, 3}

	var doubler = func(i int) int {
		return i * 2
	}

	fmt.Println(Map(testNumbers, doubler)) // 输出：[2, 4, 6]
}

func TestReduceDemo(t *testing.T) {
	// 应用实例：对数列进行求和
	var testNumbers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	var accumulator = func(previous int, current int) int {
		return current + previous
	}

	fmt.Println(Reduce(testNumbers, accumulator, 0)) // 输出：45
}

type User struct {
	Id   string
	Name string
	Age  int
	// ...
}

type Users []*User

func (users Users) Filter(filter func(*User) bool) Users {
	result := make(Users, 0, len(users))
	for _, user := range users {
		if filter(user) {
			result = append(result, user)
		}
	}
	return result
}

type Config struct {
	Id  string
	Key string
	// ...
}

type Configs []*Config

func (configs Configs) Filter(filter func(*Config) bool) Configs {
	result := make(Configs, 0, len(configs))
	for _, config := range configs {
		if filter(config) {
			result = append(result, config)
		}
	}
	return result
}
