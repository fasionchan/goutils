/*
 * Author: fasion
 * Created time: 2023-12-20 08:58:01
 * Last Modified by: fasion
 * Last Modified time: 2023-12-20 08:59:30
 */

package baseutils

import (
	"testing"
	"time"
)

func TestTimeNow(t *testing.T) {
	for i := 0; i < 100000000; i++ {
		time.Now()
	}
}
