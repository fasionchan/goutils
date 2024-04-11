/*
 * Author: fasion
 * Created time: 2023-03-24 11:59:31
 * Last Modified by: fasion
 * Last Modified time: 2024-04-11 15:35:23
 */

package baseutils

import (
	"reflect"
	"time"

	"github.com/fasionchan/goutils/stl"
)

func AlignNextTime(base time.Time, interval time.Duration, offset time.Duration) time.Time {
	result := base.Truncate(interval).Add(offset)
	for result.Before(base) {
		result = result.Add(interval)
	}
	return result
}

func NewTimes(times ...time.Time) []time.Time {
	return times
}

func MinTime(times ...time.Time) time.Time {
	return stl.Headmost(times, time.Time.Before)
}

func MaxTime(times ...time.Time) time.Time {
	return stl.Headmost(times, time.Time.After)
}

func SortTimes(times ...time.Time) []time.Time {
	return stl.Sort(times, time.Time.Before)
}

var ReflectTimeType = stl.ReflectType[time.Time]()

func WrapTimeFields(data any, wrapper func(time.Time) time.Time) error {
	return WrapTimeFieldsByReflectValue(reflect.ValueOf(data), wrapper)
}

func WrapTimeFieldsByReflectValue(value reflect.Value, wrapper func(time.Time) time.Time) error {
	valueType := value.Type()
	if valueType == ReflectTimeType {
		if value.CanSet() {
			t, ok := value.Interface().(time.Time)
			if ok {
				value.Set(reflect.ValueOf(wrapper(t)))
			}
		}

		return nil
	}

	switch value.Type().Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return nil
		}

		return WrapTimeFieldsByReflectValue(value.Elem(), wrapper)
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if err := WrapTimeFieldsByReflectValue(value.Field(i), wrapper); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if err := WrapTimeFieldsByReflectValue(value.Index(i), wrapper); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			mapValue := value.MapIndex(key)
			if mapValue.Type() == ReflectTimeType {
				t, ok := mapValue.Interface().(time.Time)
				if ok {
					value.SetMapIndex(key, reflect.ValueOf(wrapper(t)))
				}

				return nil
			}

			if err := WrapTimeFieldsByReflectValue(mapValue, wrapper); err != nil {
				return err
			}
		}
	}

	return nil
}
