package _time

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/stretchr/testify/assert"
)

func TestParseTimestamp(t *testing.T) {
	now := time.Now()

	tests := []struct {
		ts     any
		format string
		want   time.Time
	}{
		{1659312000, ParseFormatTimeUnix, time.Unix(1659312000, 0)},
		{1659312000000, ParseFormatTimeUnixMilli, time.UnixMilli(1659312000000)},
		{1659312000000000, ParseFormatTimeUnixMicro, time.UnixMicro(1659312000000000)},
		{1659312000000000000, ParseFormatTimeUnixNano, time.Unix(0, 1659312000000000000)},

		{1659312000, ParseFormatTimeUnixAuto, time.Unix(1659312000, 0)},
		{1659312000000, ParseFormatTimeUnixAuto, time.UnixMilli(1659312000000)},
		{1659312000000000, ParseFormatTimeUnixAuto, time.UnixMicro(1659312000000000)},
		{1659312000000000000, ParseFormatTimeUnixAuto, time.Unix(0, 1659312000000000000)},

		{"1659312000", ParseFormatTimeUnix, time.Unix(1659312000, 0)},
		{"1659312000000", ParseFormatTimeUnixMilli, time.UnixMilli(1659312000000)},
		{"1659312000000000", ParseFormatTimeUnixMicro, time.UnixMicro(1659312000000000)},
		{"1659312000000000000", ParseFormatTimeUnixNano, time.Unix(0, 1659312000000000000)},

		{"1659312000", ParseFormatTimeUnixAuto, time.Unix(1659312000, 0)},
		{"1659312000000", ParseFormatTimeUnixAuto, time.UnixMilli(1659312000000)},
		{"1659312000000000", ParseFormatTimeUnixAuto, time.UnixMicro(1659312000000000)},
		{"1659312000000000000", ParseFormatTimeUnixAuto, time.Unix(0, 1659312000000000000)},

		{now.Unix(), ParseFormatTimeUnix, time.Unix(now.Unix(), 0)},
		{now.UnixMilli(), ParseFormatTimeUnixMilli, time.UnixMilli(now.UnixMilli())},
		{now.UnixMicro(), ParseFormatTimeUnixMicro, time.UnixMicro(now.UnixMicro())},
		{now.UnixNano(), ParseFormatTimeUnixNano, time.Unix(0, now.UnixNano())},

		{now.Unix(), ParseFormatTimeUnixAuto, time.Unix(now.Unix(), 0)},
		{now.UnixMilli(), ParseFormatTimeUnixAuto, time.UnixMilli(now.UnixMilli())},
		{now.UnixMicro(), ParseFormatTimeUnixAuto, time.UnixMicro(now.UnixMicro())},
		{now.UnixNano(), ParseFormatTimeUnixAuto, time.Unix(0, now.UnixNano())},

		{strconv.FormatInt(now.Unix(), 10), ParseFormatTimeUnix, time.Unix(now.Unix(), 0)},
		{strconv.FormatInt(now.UnixMilli(), 10), ParseFormatTimeUnixMilli, time.UnixMilli(now.UnixMilli())},
		{strconv.FormatInt(now.UnixMicro(), 10), ParseFormatTimeUnixMicro, time.UnixMicro(now.UnixMicro())},
		{strconv.FormatInt(now.UnixNano(), 10), ParseFormatTimeUnixNano, time.Unix(0, now.UnixNano())},

		{strconv.FormatInt(now.Unix(), 10), ParseFormatTimeUnixAuto, time.Unix(now.Unix(), 0)},
		{strconv.FormatInt(now.UnixMilli(), 10), ParseFormatTimeUnixAuto, time.UnixMilli(now.UnixMilli())},
		{strconv.FormatInt(now.UnixMicro(), 10), ParseFormatTimeUnixAuto, time.UnixMicro(now.UnixMicro())},
		{strconv.FormatInt(now.UnixNano(), 10), ParseFormatTimeUnixAuto, time.Unix(0, now.UnixNano())},
	}
	for _, test := range tests {
		got, err := ParseTimestamp(test.ts, test.format)
		if err != nil {
			t.Errorf("ParseTimestamp(%v, %s) = %v", test.ts, test.format, err)
		}
		if got != test.want {
			t.Errorf("ParseTimestamp(%v, %s) = %v, want %v", test.ts, test.format, got, test.want)
		}
	}
}

type TruncateLocalTest struct {
	t time.Time
	d time.Duration
	want time.Time
}

func (test TruncateLocalTest) GetName() string {
	return fmt.Sprintf("TruncateLocal(%s, %s)", test.t.Format(time.RFC3339), test.d.String())
}

func (test TruncateLocalTest) Run(t *testing.T) {
	got := TruncateLocal(test.t, test.d)
	assert.Equal(t, test.want, got)
}

func TestTruncateLocal(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []TruncateLocalTest{
		{t: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC), d: time.Hour, want: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)},
		{t: time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC), d: time.Hour * 24, want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},

		{t: time.Date(2021, 1, 1, 8, 0, 0, 0, time.Local), d: time.Hour, want: time.Date(2021, 1, 1, 8, 0, 0, 0, time.Local)},
		{t: time.Date(2021, 1, 1, 7, 0, 0, 0, time.Local), d: time.Hour * 24, want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)},
		{t: time.Date(2021, 1, 1, 7, 0, 0, 0, time.Local), d: time.Hour * 3, want: time.Date(2021, 1, 1, 6, 0, 0, 0, time.Local)},
		{t: time.Date(2021, 1, 1, 8, 0, 0, 0, time.Local), d: time.Hour * 24, want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)},
		{t: time.Date(2021, 1, 1, 9, 0, 0, 0, time.Local), d: time.Hour * 24, want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)},
	})
}