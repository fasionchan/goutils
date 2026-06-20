package gonumplot

import (
	"strings"
	"testing"
	"time"
)

func TestChooseMajorStep(t *testing.T) {
	tests := []struct {
		span float64
		want time.Duration
	}{
		{3600, 30 * time.Minute},
		{6 * 3600, time.Hour},
		{24 * 3600, 3 * time.Hour},
		{5 * 24 * 3600, 6 * time.Hour},
	}
	for _, tt := range tests {
		if got := chooseMajorStep(tt.span); got != tt.want {
			t.Errorf("chooseMajorStep(%v) = %v, want %v", tt.span, got, tt.want)
		}
	}
}

func TestFormatDatetimeLabel(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}

	midnight := time.Date(2024, 6, 19, 0, 0, 0, 0, loc)
	if got := formatDatetimeLabel(midnight, time.Hour); got != "06-19\n00:00" {
		t.Errorf("midnight label = %q, want 06-19\\n00:00", got)
	}

	noon := time.Date(2024, 6, 19, 12, 0, 0, 0, loc)
	if got := formatDatetimeLabel(noon, time.Hour); got != "12:00" {
		t.Errorf("noon label = %q, want 12:00", got)
	}

	daily := time.Date(2024, 6, 19, 0, 0, 0, 0, loc)
	if got := formatDatetimeLabel(daily, 24*time.Hour); got != "06-19" {
		t.Errorf("daily label = %q, want 06-19", got)
	}
}

func TestDatetimeTickerTicks(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}

	// 18 小时跨度：应出现 0 点日期标签与若干整点标签
	start := time.Date(2024, 6, 19, 0, 0, 0, 0, loc).Unix()
	end := time.Date(2024, 6, 19, 18, 0, 0, 0, loc).Unix()

	ticker := newDatetimeTicker(loc)
	ticks := ticker.Ticks(float64(start), float64(end))
	if len(ticks) == 0 {
		t.Fatal("expected ticks")
	}

	var hasDate, hasTime bool
	for _, tick := range ticks {
		if strings.Contains(tick.Label, "06-19") {
			hasDate = true
		}
		if tick.Label == "12:00" {
			hasTime = true
		}
		if strings.Count(tick.Label, "06-19") > 0 && tick.Label != "06-19\n00:00" && tick.Label != "06-19" {
			t.Errorf("unexpected date in label %q", tick.Label)
		}
	}
	if !hasDate {
		t.Error("expected midnight date label")
	}
	if !hasTime {
		t.Error("expected 12:00 time label")
	}
}

func TestChartLocation(t *testing.T) {
	loc, err := chartLocation(&Chart{})
	if err != nil || loc != time.UTC {
		t.Errorf("empty timeZone: got %v err %v", loc, err)
	}

	loc, err = chartLocation(&Chart{TimeZone: "Asia/Shanghai"})
	if err != nil || loc.String() != "Asia/Shanghai" {
		t.Errorf("Asia/Shanghai: got %v err %v", loc, err)
	}

	_, err = chartLocation(&Chart{TimeZone: "Invalid/Zone"})
	if err == nil {
		t.Error("expected error for invalid timeZone")
	}
}
