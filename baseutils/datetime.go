/*
 * Author: fasion
 * Created time: 2023-03-24 11:59:31
 * Last Modified by: fasion
 * Last Modified time: 2023-03-24 11:59:44
 */

package baseutils

import "time"

func AlignNextTime(base time.Time, interval time.Duration, offset time.Duration) time.Time {
	result := base.Truncate(interval).Add(offset)
	for result.Before(base) {
		result = result.Add(interval)
	}
	return result
}
