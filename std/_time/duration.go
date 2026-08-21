package _time

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/rickb777/date/period"
)

type Duration time.Duration

func DurationFromNative(d time.Duration) Duration {
	return Duration(d)
}

func ParseIso8601Duration(s string) (Duration, error) {
	p, err := period.Parse(s)
	if err != nil {
		return 0, err
	}

	return Duration(p.Years())*Year +
		Duration(p.Months())*Month +
		Duration(p.Days())*Day +
		Duration(p.Hours())*Hour +
		Duration(p.Minutes())*Minute +
		Duration(p.Seconds())*Second, nil
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) Factor() FactoredDurationUnits {
	return d.FactorByUnitsX(DurationUnitYear, DurationUnitMonth, DurationUnitDay, DurationUnitHour, DurationUnitMinute, DurationUnitSecond)
}

func (d Duration) FactorByUnits(units []int) FactoredDurationUnits {
	return CommonDurationUnitsMapping.ValuesByKeys(units...).Factor(d)
}

func (d Duration) FactorByUnitsX(units ...int) FactoredDurationUnits {
	return d.FactorByUnits(units)
}

func (d Duration) LocaleString(locale string) string {
	return d.LocaleStringPro(locale, true, true, 2)
}

func (d Duration) LocaleStringPro(locale string, purgeZeroHead, purgeZeroTail bool, limit int) string {
	return d.Factor().LocaleStringPro(locale, purgeZeroHead, purgeZeroTail, limit)
}

func (d Duration) Parts() (int, int, int, int, int, int) {
	sign := 1
	negative := d < 0
	if negative {
		d = -d
		sign = -1
	}

	Y := d / Year
	d = d % Year

	M := d / Month
	d = d % Month

	D := d / Day
	d = d % Day

	h := d / Hour
	d = d % Hour

	m := d / Minute
	d = d % Minute

	s := d / Second

	return int(Y) * sign, int(M) * sign, int(D) * sign, int(h) * sign, int(m) * sign, int(s) * sign
}

func (d Duration) YearMonthDayDuration() (int, int, int, time.Duration) {
	Y, M, D, _, _, _ := d.Parts()
	return Y, M, D, (d % Day).Duration()
}

func (d Duration) AddTo(t time.Time) time.Time {
	if d < 0 {
		return (-d).SubFrom(t)
	}

	Y, M, D, _d := d.YearMonthDayDuration()

	return t.Add(_d).AddDate(Y, M, D)
}

func (d Duration) SubFrom(t time.Time) time.Time {
	if d < 0 {
		return (-d).AddTo(t)
	}

	Y, M, D, _d := d.YearMonthDayDuration()

	return t.Add(-_d).AddDate(-Y, -M, -D)
}

func (d Duration) Iso8601String() string {
	return period.New(d.Parts()).String()
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Iso8601String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseIso8601Duration(s)
	if err != nil {
		return err
	}

	*d = parsed

	return nil
}

func (d Duration) RandomBetween(other Duration) Duration {
	minInt := min(int64(d), int64(other))
	maxInt := max(int64(d), int64(other))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Duration(r.Int63n(maxInt-minInt+1) + minInt)
}

func DurationBetween(since, until time.Time) time.Duration {
	if since.IsZero() {
		return -1
	}

	if until.IsZero() {
		return -1
	}

	if since.After(until) {
		return -1
	}

	return until.Sub(since)
}

const (
	DurationUnitYear = iota
	DurationUnitMonth
	DurationUnitDay
	DurationUnitHour
	DurationUnitMinute
	DurationUnitSecond
	DurationUnitMillisecond
	DurationUnitMicrosecond
	DurationUnitNanosecond

	DurationUnitLocaleEn = "en"
	DurationUnitLocaleZh = "zh"
)

var (
	CommonDurationUnits = DurationUnits{
		{Unit: DurationUnitYear, Base: Year},
		{Unit: DurationUnitMonth, Base: Month},
		{Unit: DurationUnitDay, Base: Day},
		{Unit: DurationUnitHour, Base: Hour},
		{Unit: DurationUnitMinute, Base: Minute},
		{Unit: DurationUnitSecond, Base: Second},
		{Unit: DurationUnitMillisecond, Base: Millisecond},
		{Unit: DurationUnitMicrosecond, Base: Microsecond},
		{Unit: DurationUnitNanosecond, Base: Nanosecond},
	}

	CommonDurationUnitsMapping = CommonDurationUnits.MappingByUnit()

	DurationUnitLocales = map[string]map[int]string{
		"zh": {
			DurationUnitYear:   "年",
			DurationUnitMonth:  "月",
			DurationUnitDay:    "天",
			DurationUnitHour:   "时",
			DurationUnitMinute: "分",
			DurationUnitSecond: "秒",
		},
		"en": {
			DurationUnitYear:   "yr",
			DurationUnitMonth:  "mo",
			DurationUnitDay:    "d",
			DurationUnitHour:   "hr",
			DurationUnitMinute: "min",
			DurationUnitSecond: "sec",
		},
	}
)

type DurationUnit struct {
	Unit int
	Base Duration
}

func (unit DurationUnit) GetUnit() int {
	return unit.Unit
}

func (unit DurationUnit) factor(ptr *Duration) (result FactoredDurationUnit) {
	result.DurationUnit = unit

	if ptr == nil {
		return
	}

	d := *ptr
	if d < unit.Base {
		return
	}

	result.Factor = d / unit.Base
	*ptr = d % unit.Base

	return
}

type DurationUnits []DurationUnit

func (uints DurationUnits) Factor(d Duration) FactoredDurationUnits {
	return stl.MapUnary(uints, DurationUnit.factor, &d)
}

func (units DurationUnits) MappingByUnit() DurationUnitMappingByInt {
	return stl.MappingByKey(units, DurationUnit.GetUnit)
}

type DurationUnitMappingByInt map[int]DurationUnit

func (mapping DurationUnitMappingByInt) ValuesByKeys(keys ...int) DurationUnits {
	return stl.MapValuesByKeys(mapping, keys...)
}

type FactoredDurationUnit struct {
	DurationUnit
	Factor Duration
}

func (unit FactoredDurationUnit) IsFactorZero() bool {
	return unit.Factor == 0
}

func (unit FactoredDurationUnit) LocaleString(localeName string) string {
	locale := DurationUnitLocales[localeName]
	if locale == nil {
		return ""
	}

	return fmt.Sprintf("%d%s", int64(unit.Factor), locale[unit.Unit])
}

type FactoredDurationUnits []FactoredDurationUnit

func (units FactoredDurationUnits) PurgeZero() FactoredDurationUnits {
	return stl.Purge(units, FactoredDurationUnit.IsFactorZero)
}

func (units FactoredDurationUnits) PurgeZeroHead() FactoredDurationUnits {
	return stl.PurgeHead(units, FactoredDurationUnit.IsFactorZero)
}

func (units FactoredDurationUnits) PurgeZeroTail() FactoredDurationUnits {
	return stl.PurgeTail(units, FactoredDurationUnit.IsFactorZero)
}

func (units FactoredDurationUnits) LocaleStringParts(locale string) types.Strings {
	return stl.MapUnary(units, FactoredDurationUnit.LocaleString, locale)
}

func (units FactoredDurationUnits) LocaleString(locale string) string {
	return units.LocaleStringParts(locale).Join("")
}

func (units FactoredDurationUnits) LocaleStringPro(locale string, purgeZeroHead, purgeZeroTail bool, limit int) string {
	if purgeZeroHead {
		units = units.PurgeZeroHead()
	}

	units = units.Limit(limit)

	if purgeZeroTail {
		units = units.PurgeZeroTail()
	}

	return units.LocaleString(locale)
}

func (units FactoredDurationUnits) Limit(limit int) FactoredDurationUnits {
	n := len(units)
	if limit < 0 {
		limit = 0
	} else if limit > n {
		limit = n
	}

	return units[:limit]
}

var (
	DurationLocaleString    = Duration.LocaleString
	DurationLocaleStringPro = Duration.LocaleStringPro
)

func NativeDurationLocaleString(d time.Duration, locale string) string {
	return DurationFromNative(d).LocaleString(locale)
}

func NativeDurationLocaleStringPro(d time.Duration, locale string, purgeZeroHead, purgeZeroTail bool, limit int) string {
	return DurationFromNative(d).LocaleStringPro(locale, purgeZeroHead, purgeZeroTail, limit)
}
