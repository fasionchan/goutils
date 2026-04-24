package timex

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/rickb777/date/period"
)

type Duration time.Duration

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
