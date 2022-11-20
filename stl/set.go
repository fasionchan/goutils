/*
 * Author: fasion
 * Created time: 2022-11-19 18:46:55
 * Last Modified by: fasion
 * Last Modified time: 2022-11-20 12:03:17
 */

package stl

type Set[Data comparable] map[Data]struct{}

func NewSet[Data comparable](datas ...Data) Set[Data] {
	return NewEmptySet[Data]().PushX(datas...)
}

func NewEmptySet[Data comparable]() Set[Data] {
	return Set[Data]{}
}

func (set Set[Data]) Contain(data Data) bool {
	_, ok := set[data]
	return ok
}

func (set Set[Data]) Len() int {
	return len(set)
}

func (set Set[Data]) Empty() bool {
	return set.Len() == 0
}

func (set Set[Data]) Equal(other Set[Data]) bool {
	if set.Len() != other.Len() {
		return false
	}

	for data := range set {
		if !other.Contain(data) {
			return false
		}
	}

	return true
}

func (set Set[Data]) Slice() []Data {
	return MapKeys(set)
}

func (set Set[Data]) Dup() Set[Data] {
	return DupMap(set)
}

func (set Set[Data]) Push(data Data) Set[Data] {
	set[data] = struct{}{}
	return set
}

func (set Set[Data]) PushX(datas ...Data) Set[Data] {
	ForEach(datas, set.Push)
	// duplicated
	// for _, data := range datas {
	// 	set.Push(data)
	// }
	// return set
}

func (set Set[Data]) Pop(data Data) Set[Data] {
	delete(set, data)
	return set
}

func (set Set[Data]) Add(other Set[Data]) Set[Data] {
	return ConcatMapInplace(set, other)
}

func (set Set[Data]) Sub(other Set[Data]) Set[Data] {
	return BatchDeleteMapFromAnother(set, other)
}

func (set Set[Data]) Union(other Set[Data]) Set[Data] {
	return ConcatMap(set, other)
}

func (set Set[Data]) Difference(other Set[Data]) Set[Data] {
	return set.Dup().Sub(other)
}

func (set Set[Data]) SymmetricDifference(other Set[Data]) Set[Data] {
	result := Set[Data]{}

	for data := range set {
		if !other.Contain(data) {
			result.Push(data)
		}
	}

	for data := range other {
		if !set.Contain(data) {
			result.Push(data)
		}
	}

	return result
}

func (set Set[Data]) Intersection(other Set[Data]) Set[Data] {
	result := Set[Data]{}
	for data := range set {
		if other.Contain(data) {
			result.Push(data)
		}
	}
	return result
}
