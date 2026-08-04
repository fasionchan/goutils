package _time

import "time"

func NextDay(d time.Time) time.Time {
	return d.AddDate(0, 0, 1)
}

func PreviousDay(d time.Time) time.Time {
	return d.AddDate(0, 0, -1)
}

func NextDaySeq(d time.Time) func (yield func(time.Time) bool) {
	return func(yield func(time.Time) bool) {
		for yield(d) {
			d = NextDay(d)
		}
	}
}

func PreviousDaySeq(d time.Time) func (yield func(time.Time) bool) {
	return func(yield func(time.Time) bool) {
		for yield(d) {
			d = PreviousDay(d)
		}
	}
}

func NextDays(d time.Time, n int) []time.Time {
	days := make([]time.Time, 0, n)
	for i := 0; i < n; i++ {
		days = append(days, NextDay(d))
		d = NextDay(d)
	}
	return days
}

func PreviousDays(d time.Time, n int) []time.Time {
	days := make([]time.Time, 0, n)
	for i := 0; i < n; i++ {
		days = append(days, PreviousDay(d))
		d = PreviousDay(d)
	}
	return days
}