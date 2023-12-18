/*
 * Author: fasion
 * Created time: 2023-03-24 11:59:31
 * Last Modified by: fasion
 * Last Modified time: 2023-12-18 09:59:38
 */

package baseutils

import (
	"time"

	"github.com/fasionchan/goutils/stl"
)

func AlignNextTime(base time.Time, interval time.Duration, offset time.Duration) time.Time {
	result := base.Truncate(interval).Add(offset)
	for result.Before(base) {
		result = result.Add(interval)
	}
	return result
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
