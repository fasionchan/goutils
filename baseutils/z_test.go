/*
 * Author: fasion
 * Created time: 2023-04-20 15:22:26
 * Last Modified by: fasion
 * Last Modified time: 2023-12-06 11:03:35
 */

package baseutils

import (
	"fmt"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {

}

func TestTimeFromUnit(t *testing.T) {
	var ts int64 = 1701419596802
	fmt.Println(time.Unix(ts/1000, ts%1000*1000000))
	fmt.Println(time.Unix(0, ts*1000000))
	fmt.Println(time.Unix(0, ts*10000000))
}
