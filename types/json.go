package types

import (
	"github.com/fasionchan/goutils/stl"
)

type JsonAny = interface{}
type JsonObject = map[string]JsonAny
type JsonMap = JsonObject
type JsonArray = []JsonAny

func NewJsonArray(datas ...JsonAny) JsonArray {
	return datas
}

type SmartJsonMap JsonMap

func NewSmartJsonMap() SmartJsonMap {
	return SmartJsonMap{}
}

func (m SmartJsonMap) Len() int {
	return len(m)
}

func (m SmartJsonMap) Empty() bool {
	return m.Len() == 0
}

func (m SmartJsonMap) Dup() SmartJsonMap {
	return stl.DupMap(m)
}

func (m SmartJsonMap) Contains(key string) (ok bool) {
	_, ok = m[key]
	return
}

func (m SmartJsonMap) With(key string, value any) SmartJsonMap {
	m[key] = value
	return m
}

func (m SmartJsonMap) Merge(other SmartJsonMap) SmartJsonMap {
	return stl.ConcatMapInplace(m, other)
}

func (m SmartJsonMap) BatchMerge(others ...SmartJsonMap) SmartJsonMap {
	for _, other := range others {
		m.Merge(other)
	}
	return m
}

func (m SmartJsonMap) Self() SmartJsonMap {
	return m
}

func (m SmartJsonMap) Keys() Strings {
	return stl.MapKeys(m)
}

func (m SmartJsonMap) ValuesByKeys(keys []string) []JsonAny {
	return stl.MapValuesByKeys(m, keys...)
}

type SmartJsonMaps = []SmartJsonMap
