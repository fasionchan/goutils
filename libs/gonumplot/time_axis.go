package gonumplot

import (
	"fmt"
	"time"

	"gonum.org/v1/plot"
)

// datetimeTicker 将 Unix 秒时间戳格式化为可读的日期时间刻度。
type datetimeTicker struct {
	loc *time.Location
}

func newDatetimeTicker(loc *time.Location) datetimeTicker {
	if loc == nil {
		loc = time.UTC
	}
	return datetimeTicker{loc: loc}
}

func (dt datetimeTicker) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		return nil
	}

	span := max - min
	step := chooseMajorStep(span)
	for span/float64(step) > 8 {
		step = step * 2
		if step > 24*time.Hour {
			break
		}
	}

	start := alignDown(time.Unix(int64(min), 0).In(dt.loc), step)
	end := time.Unix(int64(max), 0).In(dt.loc)

	var ticks []plot.Tick
	for t := start; !t.After(end); t = t.Add(step) {
		v := float64(t.Unix())
		if v < min || v > max {
			continue
		}
		ticks = append(ticks, plot.Tick{
			Value: v,
			Label: formatDatetimeLabel(t, step),
		})
	}
	return ticks
}

// chooseMajorStep 按时间跨度选择初始主刻度间隔。
func chooseMajorStep(spanSeconds float64) time.Duration {
	switch {
	case spanSeconds <= 2*3600:
		return 30 * time.Minute
	case spanSeconds <= 12*3600:
		return time.Hour
	case spanSeconds <= 2*24*3600:
		return 3 * time.Hour
	case spanSeconds <= 7*24*3600:
		return 6 * time.Hour
	case spanSeconds <= 30*24*3600:
		return 12 * time.Hour
	default:
		return 24 * time.Hour
	}
}

// alignDown 将时间向下对齐到 step 边界（在 loc 时区下）。
func alignDown(t time.Time, step time.Duration) time.Time {
	loc := t.Location()
	switch {
	case step >= 24*time.Hour:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	case step >= time.Hour:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
	case step >= 30*time.Minute:
		minute := t.Minute() / 30 * 30
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, loc)
	default:
		minute := t.Minute() / 10 * 10
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, loc)
	}
}

// formatDatetimeLabel 按常规监控图表策略格式化刻度标签：
// - 0 点显示日期（跨天边界）
// - 整点显示 HH:MM
// - 其余刻度仅显示时间，不重复日期
func formatDatetimeLabel(t time.Time, step time.Duration) string {
	if t.Hour() == 0 && t.Minute() == 0 {
		if step >= 24*time.Hour {
			return t.Format("01-02")
		}
		return t.Format("01-02\n15:04")
	}
	if t.Minute() == 0 {
		return t.Format("15:04")
	}
	return t.Format("15:04")
}

// chartLocation 解析 Chart 配置的 IANA 时区。
func chartLocation(chart *Chart) (*time.Location, error) {
	if chart.TimeZone == "" {
		return time.UTC, nil
	}
	loc, err := time.LoadLocation(chart.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("invalid timeZone %q: %w", chart.TimeZone, err)
	}
	return loc, nil
}
