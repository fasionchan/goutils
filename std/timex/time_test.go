package timex

import (
	"strconv"
	"testing"
	"time"
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
