/*
 * Author: fasion
 * Created time: 2023-03-24 11:59:31
 * Last Modified by: fasion
 * Last Modified time: 2025-04-28 09:05:15
 */

package baseutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fasionchan/goutils/stl"
)

const (
	DurationDay      = time.Hour * 24
	DurationMonth    = DurationDay * 30
	MaxDurationMonth = DurationDay * 31
	DurationYear     = DurationDay * 365
	MaxDurationYear  = DurationDay * 366
)

func AlignNextTime(base time.Time, interval time.Duration, offset time.Duration) time.Time {
	result := base.Truncate(interval).Add(offset)
	for result.Before(base) {
		result = result.Add(interval)
	}
	return result
}

var year2050 = time.Date(2050, 1, 1, 0, 0, 0, 0, time.Now().Local().Location())
var year3000 = time.Date(3000, 1, 1, 0, 0, 0, 0, time.Now().Local().Location())
var year5000 = time.Date(5000, 1, 1, 0, 0, 0, 0, time.Now().Local().Location())

func GetYear2050() time.Time {
	return year2050
}

func GetYear3000() time.Time {
	return year3000
}

func GetYear5000() time.Time {
	return year5000
}

func DayOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Today() time.Time {
	return DayOf(time.Now().Local())
}

func Tomorrow() time.Time {
	return Today().Add(time.Hour * 24)
}

func Yesterday() time.Time {
	return Today().Add(-time.Hour * 24)
}

type IntraDayTime time.Duration

func ParseIntraDayTime(s string) (IntraDayTime, error) {
	var hours, minutes, seconds, nanoseconds time.Duration

	var parts = strings.Split(s, ".")
	if _, err := fmt.Sscanf(parts[0], "%d:%d:%d", &hours, &minutes, &seconds); err != nil {
		return 0, err
	}

	// parse nanosecond if any
	if len(parts) > 1 {
		if ns := parts[1]; ns != "" {
			if _, err := fmt.Sscanf(ns, "%d", &nanoseconds); err != nil {
				return 0, err
			}
		}
	}

	return IntraDayTime(hours*time.Hour + minutes*time.Minute + seconds*time.Second + nanoseconds*time.Nanosecond), nil
}

func MustParseIntraDayTime(s string) IntraDayTime {
	t, err := ParseIntraDayTime(s)
	if err != nil {
		panic(err)
	}

	return t
}

func ParseFormattedIntraDayTime(layout, value string) (IntraDayTime, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return 0, err
	}

	r := IntraDayTime(time.Hour)*IntraDayTime(t.Hour()) +
		IntraDayTime(time.Minute)*IntraDayTime(t.Minute()) +
		IntraDayTime(time.Second)*IntraDayTime(t.Second()) +
		IntraDayTime(time.Nanosecond)*IntraDayTime(t.Nanosecond())

	return r, nil
}

func MustParseFormattedIntraDayTime(layout, value string) IntraDayTime {
	t, err := ParseFormattedIntraDayTime(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}

func (t IntraDayTime) Duration() time.Duration {
	return time.Duration(t)
}

func (t IntraDayTime) String() string {
	hours, minutes, seconds, nanoseconds := t.Parts()
	return fmt.Sprintf("%02d:%02d:%02d.%09d", hours, minutes, seconds, nanoseconds)
}

func (t IntraDayTime) Format(layout string) string {
	hours, minutes, seconds, nanoseconds := t.Parts()
	return time.Date(0, 0, 0, hours, minutes, seconds, nanoseconds, time.Local).Format(layout)
}

func (t IntraDayTime) Parts() (int, int, int, int) {
	d := time.Duration(t)

	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second
	d -= seconds * time.Second

	return int(hours), int(minutes), int(seconds), int(d)
}

func (t IntraDayTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *IntraDayTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseIntraDayTime(s)
	if err != nil {
		return err
	}

	*t = parsed

	return nil
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

type DateTime time.Time

func (t DateTime) Native() time.Time {
	return time.Time(t)
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.DateTime))
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	_time, err := time.ParseInLocation(time.DateTime, str, time.Local)
	if err != nil {
		return err
	}

	*t = DateTime(_time)

	return nil
}
