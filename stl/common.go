/*
 * Author: fasion
 * Created time: 2023-03-24 08:46:11
 * Last Modified by: fasion
 * Last Modified time: 2023-03-24 08:46:20
 */

package stl

func Dup[Data any, Ptr ~*Data](ptr Ptr) Ptr {
	dup := *ptr
	return &dup
}
