/*
 * Author: fasion
 * Created time: 2022-11-19 17:56:20
 * Last Modified by: fasion
 * Last Modified time: 2022-11-20 11:57:42
 */

package stl

import "strings"

type StringSlice = []string

type Strings []string

func NewStrings(strs ...string) Strings {
	return Strings(strs)
}

func (strs Strings) Len() int {
	return len(strs)
}

func (strs Strings) Empty() bool {
	return strs.Len() == 0
}

func (strs Strings) IndexOf(str string) int {
	return IndexOf(strs, str)
}

func (strs Strings) StringSlice() StringSlice {
	return StringSlice(strs)
}

func (strs Strings) StringSet() StringSet {
	return NewSet(strs...)
}

func (strs Strings) Dup() Strings {
	return DupSlice(strs)
}

func (strs Strings) Append(more ...string) Strings {
	return append(strs, more...)
}

func (strs Strings) Concat(others ...Strings) Strings {
	return ConcatSlicesTo(strs, others...)
}

func (strs Strings) Sort() Strings {
	return Sort(strs, func(a, b string) bool { return a < b })
}

func (strs Strings) AnyMatch(test func(string) bool) bool {
	return AnyMatch(strs, test)
}

func (strs Strings) AllMatch(test func(string) bool) bool {
	return AllMatch(strs, test)
}

func (strs Strings) ForEach(handler func(string)) {
	ForEach(strs, handler)
}

func (strs Strings) ForEachCustom(handler func(int, string, Strings)) {
	ForEachCustom(strs, handler)
}

func (strs Strings) Filter(filter func(string) bool) Strings {
	return Filter(strs, filter)
}

func (strs Strings) Purge(filter func(string) bool) Strings {
	return Purge(strs, filter)
}

func (strs Strings) PurgeZero() Strings {
	return PurgeZero(strs)
}

func (strs Strings) Map(mapper func(string) string) Strings {
	return Map(strs, mapper)
}

func (strs Strings) TrimSpace() Strings {
	return strs.Map(strings.TrimSpace)
}

// go example
func (strs Strings) Split(seps ...string) Strings {
	return Reduce(seps, func(sep string, result Strings) Strings {
		return MapAndConcat(result, func(str string) Strings {
			return strings.Split(str, sep)
		})
	}, strs)
}

func (strs Strings) Unique() Strings {
	return UniqueFast(strs)
}

func (strs Strings) Join(sep string) string {
	return strings.Join(strs.StringSlice(), sep)
}

type StringSet = Set[string]
