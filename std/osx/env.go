package osx

import (
	"fmt"
	"os"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

type EnvironMap map[string]string

func GetProcessEnviron() EnvironMap {
	return stl.BuildMap[EnvironMap](os.Environ(), func(item string) (string, string) {
		fields := strings.Split(item, "=")
		return fields[0], fields[1]
	})
}

func (m EnvironMap) Keys() []string {
	return stl.MapKeys(m)
}

func (m EnvironMap) LookupEnv(key string, realfirst bool) (value string, ok bool) {
	if realfirst {
		value, ok = os.LookupEnv(key)
		if ok {
			return
		}

		value, ok = m[key]

		return
	} else {
		value, ok = m[key]
		if ok {
			return
		}

		value, ok = os.LookupEnv(key)

		return
	}
}

func (m EnvironMap) Getenv(key string, realfirst bool) (value string) {
	value, _ = m.LookupEnv(key, realfirst)
	return
}

func (m EnvironMap) GetGetter(realfirst bool) func(string) string {
	return func(key string) string {
		return m.Getenv(key, realfirst)
	}
}

func (m EnvironMap) GetLooker(realfirst bool) func(string) (string, bool) {
	return func(key string) (string, bool) {
		return m.LookupEnv(key, realfirst)
	}
}

func (m EnvironMap) Dup() EnvironMap {
	return stl.DupMap(m)
}

func (m EnvironMap) With(name, value string) EnvironMap {
	m[name] = value
	return m
}

func (m EnvironMap) Concat(all EnvironMap) EnvironMap {
	return stl.ConcatMapInplace(m, all)
}

func (m EnvironMap) Format() types.Strings {
	return stl.MapMapToSlice[types.Strings](m, func(name, value string, _ EnvironMap) string {
		return fmt.Sprintf("%s=%s", name, value)
	})
}