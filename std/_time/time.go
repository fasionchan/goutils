package _time

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"

	ParseFormatTimeUnix      = "Unix"
	ParseFormatTimeUnixMilli = "UnixMilli"
	ParseFormatTimeUnixMicro = "UnixMicro"
	ParseFormatTimeUnixNano  = "UnixNano"
	ParseFormatTimeUnixAuto  = "UnixAuto"

	MixTimestampNano  = 1000000000000000000
	MixTimestampMicro = 1000000000000000
	MixTimestampMilli = 1000000000000
	MixTimestamp      = 1000000000
)

var (
	TimeAligners = map[string]func(time.Time) time.Time{
		"Year":       AlignTimeByYear,
		"YearLocal":  AlignTimeByYearLocal,
		"Month":      AlignTimeByMonth,
		"MonthLocal": AlignTimeByMonthLocal,
		"Week":       AlignTimeByWeek,
		"Day":        AlignTimeByDay,
		"Hour":       AlignTimeByHour,
		"Minute":     AlignTimeByMinute,
	}

	FormatRfc3339     = NewTimeFormatter(time.RFC3339)
	FormatRfc3339Nano = NewTimeFormatter(time.RFC3339Nano)
	FormatIsoTime     = FormatRfc3339

	ParseRfc3339     = NewTimeParser(time.RFC3339)
	ParseRfc3339Nano = NewTimeParser(time.RFC3339Nano)
	ParseIsoTime     = ParseRfc3339
)

func AlignTimeByYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func AlignTimeByYearWithLocation(t time.Time, location *time.Location) time.Time {
	if location != nil {
		t = t.In(location)
	}
	return AlignTimeByYear(t)
}

func AlignTimeByYearLocal(t time.Time) time.Time {
	return AlignTimeByYear(t.Local())
}

func AlignTimeByQuarter(t time.Time) time.Time {
	return time.Date(t.Year(), (t.Month()-1)/3*3+1, 1, 0, 0, 0, 0, t.Location())
}

func AlignTimeByQuarterWithLocation(t time.Time, location *time.Location) time.Time {
	if location != nil {
		t = t.In(location)
	}
	return AlignTimeByQuarter(t)
}

func AlignTimeByQuarterLocal(t time.Time) time.Time {
	return AlignTimeByQuarter(t.Local())
}

func AlignTimeByMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func AlignTimeByMonthWithLocation(t time.Time, location *time.Location) time.Time {
	if location != nil {
		t = t.In(location)
	}
	return AlignTimeByMonth(t)
}

func AlignTimeByMonthLocal(t time.Time) time.Time {
	return AlignTimeByMonth(t.Local())
}

func AlignTimeByWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
}

func AlignTimeByWeekWithLocation(t time.Time, location *time.Location) time.Time {
	if location != nil {
		t = t.In(location)
	}
	return AlignTimeByWeek(t)
}

func AlignTimeByWeekLocal(t time.Time) time.Time {
	return AlignTimeByWeek(t.Local())
}

func AlignTimeByDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func AlignTimeByDayWithLocation(t time.Time, location *time.Location) time.Time {
	if location != nil {
		t = t.In(location)
	}
	return AlignTimeByDay(t)
}

func AlignTimeByDayLocal(t time.Time) time.Time {
	return AlignTimeByDay(t.Local())
}

func AlignTimeByHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

func AlignTimeByMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

func FormatTime(t time.Time, fmt string, zeroPlaceHolder string) string {
	if t.IsZero() {
		return zeroPlaceHolder
	}

	if fmt == "" {
		fmt = DefaultTimeFormat
	}

	return t.Format(fmt)
}

func NewTimeFormatter(fmt string) func(t time.Time, zeroPlaceHolder string) string {
	return func(t time.Time, zeroPlaceHolder string) string {
		return FormatTime(t, fmt, zeroPlaceHolder)
	}
}

func NewTimeParser(layout string) func(s string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		return time.Parse(layout, s)
	}
}

func ParseTimestampAuto(ts any) (time.Time, error) {
	return ParseTimestamp(ts, ParseFormatTimeUnixAuto)
}

func ParseTimestampSecond(ts any) (time.Time, error) {
	return ParseTimestamp(ts, ParseFormatTimeUnix)
}

func ParseTimestampMilliSecond(ts any) (time.Time, error) {
	return ParseTimestamp(ts, ParseFormatTimeUnixMilli)
}

func ParseTimestampMicroSecond(ts any) (time.Time, error) {
	return ParseTimestamp(ts, ParseFormatTimeUnixMicro)
}

func ParseTimestampNanoSecond(ts any) (time.Time, error) {
	return ParseTimestamp(ts, ParseFormatTimeUnixNano)
}

func ParseTimestamp(ts any, format string) (time.Time, error) {
	v := reflect.ValueOf(ts)
	var value int64
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		value = int64(v.Float())
	case reflect.String:
		return ParseTimestampStr(ts.(string), format)
	default:
		return time.Time{}, fmt.Errorf("ParseTimestamp[%v]", v.Interface())
	}

	return ParseTimestampInt(value, format), nil
}

func ParseTimestampStr(tsStr, format string) (time.Time, error) {
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("ParseTimestampStrFailed[%d], err:%w", ts, err)
	}

	return ParseTimestampInt(ts, format), nil
}

func ParseTimestampStrAuto(tsStr string) (time.Time, error) {
	if n := len(tsStr); n >= 19 {
		return ParseTimestampStr(tsStr, ParseFormatTimeUnixNano)
	} else if n >= 16 {
		return ParseTimestampStr(tsStr, ParseFormatTimeUnixMicro)
	} else if n >= 13 {
		return ParseTimestampStr(tsStr, ParseFormatTimeUnixMilli)
	} else {
		return ParseTimestampStr(tsStr, ParseFormatTimeUnix)
	}
}

func ParseTimestampInt(ts int64, format string) time.Time {
	switch format {
	case ParseFormatTimeUnix:
		return time.Unix(ts, 0)
	case ParseFormatTimeUnixMilli:
		return time.UnixMilli(ts)
	case ParseFormatTimeUnixMicro:
		return time.UnixMicro(ts)
	case ParseFormatTimeUnixNano:
		return time.Unix(0, ts)
	case ParseFormatTimeUnixAuto:
		return ParseTimestampIntAuto(ts)
	default:
		return time.Time{}
	}
}

func ParseTimestampIntAuto(ts int64) time.Time {
	if ts >= MixTimestampNano {
		return time.Unix(0, ts)
	} else if ts >= MixTimestampMicro {
		return time.UnixMicro(ts)
	} else if ts >= MixTimestampMilli {
		return time.UnixMilli(ts)
	} else {
		return time.Unix(ts, 0)
	}
}
