/*
 * Author: fasion
 * Created time: 2022-11-19 17:56:20
 * Last Modified by: fasion
 * Last Modified time: 2025-09-23 13:45:49
 */

package types

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"strings"

	"github.com/fasionchan/goutils/stl"
)

type String string

func (s String) IsPrefixOf(str string) bool {
	return strings.HasPrefix(str, string(s))
}

func (s String) IsSuffixOf(str string) bool {
	return strings.HasSuffix(str, string(s))
}

type StringSet = stl.Set[string]

type StringSlice = []string

type Strings []string

func NewStrings(strs ...string) Strings {
	return Strings(strs)
}

func NewStringsFromStringSlice(strs StringSlice) Strings {
	return strs
}

func NewStringsFromSet(set StringSet) Strings {
	return NewStringsFromStringSlice(set.Slice())
}

func SplitToStrings(s, sep string) Strings {
	return strings.Split(s, sep)
}

func (strs Strings) Len() int {
	return len(strs)
}

func (strs Strings) Empty() bool {
	return strs.Len() == 0
}

func (strs Strings) Native() StringSlice {
	return StringSlice(strs)
}

func (strs Strings) ToSet() StringSet {
	return stl.NewSet(strs...)
}

func (strs Strings) Count() stl.Counter[string] {
	return stl.Count(strs)
}

func (strs Strings) Dup() Strings {
	return stl.DupSlice(strs)
}

func (strs Strings) EnsureNotNil() Strings {
	if strs == nil {
		return Strings{}
	}
	return strs
}

func (strs Strings) IndexOf(str string) int {
	return stl.IndexOf(strs, str)
}

func (strs Strings) Append(more ...string) Strings {
	return append(strs, more...)
}

func (strs Strings) Concat(others ...Strings) Strings {
	return stl.ConcatSlicesTo(strs, others...)
}

func (strs Strings) InplaceShuffle() Strings {
	return stl.Shuffle(strs, nil)
}

func (strs Strings) InplaceShufflePro(rand_ *rand.Rand) Strings {
	return stl.Shuffle(strs, rand_)
}

func (strs Strings) InplaceSort() Strings {
	return stl.Sort(strs, StringComparer)
}

func (strs Strings) Sort() Strings {
	return stl.Sort(strs.Dup(), StringComparer)
}

func (strs Strings) AnyHasPrefix(prefix string) bool {
	return strs.AnyMatch(String(prefix).IsPrefixOf)
}

func (strs Strings) AnyHasSuffix(suffix string) bool {
	return strs.AnyMatch(String(suffix).IsSuffixOf)
}

func (strs Strings) AnyMatch(test func(string) bool) bool {
	return stl.AnyMatch(strs, test)
}

func (strs Strings) AllMatch(test func(string) bool) bool {
	return stl.AllMatch(strs, test)
}

func (strs Strings) FirstOneOrZero() string {
	return stl.FirstOneOrZero(strs)
}

func (strs Strings) ForEach(handler func(string)) {
	stl.ForEach(strs, handler)
}

func (strs Strings) ForEachPro(handler func(int, string, Strings)) {
	stl.ForEachPro(strs, handler)
}

func (strs Strings) Filter(filter func(string) bool) Strings {
	return stl.Filter(strs, filter)
}

func (strs Strings) Purge(filter func(string) bool) Strings {
	return stl.Purge(strs, filter)
}

func (strs Strings) PurgeZero() Strings {
	return stl.PurgeZero(strs)
}

func (strs Strings) Map(mapper func(string) string) Strings {
	return stl.Map(strs, mapper)
}

func (strs Strings) MapWithSprintf(format string) Strings {
	return strs.Map(func(s string) string {
		return fmt.Sprintf(format, s)
	})
}

func (strs Strings) RandomOneOrZero() string {
	return stl.RandomOneOrZero(strs, nil)
}

func (strs Strings) RandomOneOrZeroPro(rand_ *rand.Rand) string {
	return stl.RandomOneOrZero(strs, rand_)
}

func (strs Strings) ReverseInplace() Strings {
	return stl.Reverse(strs)
}

func (strs Strings) ReverseDup() Strings {
	return stl.Reverse(strs.Dup())
}

func (strs Strings) SliceOfAny() []any {
	return stl.ToTypelessSlice(strs)
}

func (strs Strings) TrimSpace() Strings {
	return strs.Map(strings.TrimSpace)
}

func (strs Strings) ToLower() Strings {
	return strs.Map(strings.ToLower)
}

func (strs Strings) ToUpper() Strings {
	return strs.Map(strings.ToUpper)
}

// go example
func (strs Strings) Split(seps ...string) Strings {
	return stl.Reduce(seps, func(result Strings, sep string) Strings {
		return stl.MapAndConcat(result, func(str string) Strings {
			return strings.Split(str, sep)
		})
	}, strs)
}

func (strs Strings) Unique() Strings {
	return strs.UniqueBySet()
}

func (strs Strings) UniqueBySet() Strings {
	return stl.UniqueBySet(strs)
}

func (strs Strings) UniqueByInplaceSort() Strings {
	return strs.InplaceSort().UniqueSorteds()
}

func (strs Strings) UniqueBySort() Strings {
	return strs.Sort().UniqueSorteds()
}

func (strs Strings) UniqueSorteds() Strings {
	return stl.UniqueSorteds(strs)
}

func (strs Strings) Join(sep string) string {
	return strings.Join(strs.Native(), sep)
}

func (strs Strings) JoinByComma() string {
	return strs.Join(",")
}

func (strs Strings) JoinByDot() string {
	return strs.Join(".")
}

func (strs Strings) JoinByEmptyLine() string {
	return strs.Join("\n\n")
}

func (strs Strings) JoinByNewLine() string {
	return strs.Join("\n")
}

func (strs Strings) JoinByMarkdownSplitLine() string {
	return strs.Join("\n---\n")
}

func (strs Strings) Equal(other Strings) bool {
	return stl.SliceEqual(strs, other)
}

func (strs Strings) Contain(str string) bool {
	return stl.ComparableSlice[string](strs).Contain(str)
}

func (strs Strings) ItemsUnique() bool {
	return strs.ToSet().Len() == strs.Len()
}

func StringComparer(a, b string) bool {
	return a < b
}

type CsvRecord = CommaSeparatedValueRecord

type CommaSeparatedValueRecord string

func (s CommaSeparatedValueRecord) Native() string {
	return string(s)
}

func (s CommaSeparatedValueRecord) Values() Strings {
	values, _ := csv.NewReader(strings.NewReader(string(s))).Read()
	return values
}

func (s CommaSeparatedValueRecord) ValidValues() Strings {
	return s.Values().PurgeZero()
}

type StringsMappingByString stl.Mapping[string, Strings]

func (mapping StringsMappingByString) KeyValuePairs() KeyStringsPairs[string] {
	return (KeyStringsPairs[string])(stl.MapKeyValuePairs(mapping))
}

type StringsMapping[Key comparable] stl.Mapping[Key, Strings]

func (mapping StringsMapping[Key]) KeyPairs() stl.KeyValuePairs[Key, Strings] {
	return mapping.KeyPairs()
}

type KeyStringPair[Key any] stl.KeyValuePair[Key, Strings]

type KeyStringsPairs[Key any] stl.KeyValuePairs[Key, Strings]

func (pairs KeyStringsPairs[Key]) Native() stl.KeyValuePairs[Key, Strings] {
	return (stl.KeyValuePairs[Key, Strings])(pairs)
}

func (pairs KeyStringsPairs[Key]) SortByValueLen() KeyStringsPairs[Key] {
	pairs.Native().SortByValue(func(a, b Strings) bool {
		return len(a) < len(b)
	})
	return pairs
}

func (pairs KeyStringsPairs[Key]) SortByValueLenDesc() KeyStringsPairs[Key] {
	pairs.Native().SortByValue(func(a, b Strings) bool {
		return len(a) > len(b)
	})
	return pairs
}
