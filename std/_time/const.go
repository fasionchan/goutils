package _time

import (
	"time"
)

const (
	DaysInMonth = 30
	DaysInYear  = 365

	Nanosecond = Duration(time.Nanosecond)
	Microsecond = Duration(time.Microsecond)
	Millisecond = Duration(time.Millisecond)
	Second = Duration(time.Second)
	Minute = Duration(time.Minute)
	Hour   = Duration(time.Hour)
	Day    = Hour * 24
	Month  = Day * DaysInMonth
	Year   = Day * DaysInYear
)
