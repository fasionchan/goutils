/*
 * Author: fasion
 * Created time: 2023-12-20 08:58:01
 * Last Modified by: fasion
 * Last Modified time: 2024-04-11 15:24:31
 */

package baseutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNow(t *testing.T) {
	for i := 0; i < 100000000; i++ {
		time.Now()
	}
}

func TestWrapTimeFields(t *testing.T) {
	timeValue := time.Now().UTC()
	utcHour := timeValue.Hour()
	localHour := timeValue.Local().Hour()

	if err := WrapTimeFields(timeValue, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, timeValue.Hour(), utcHour)

	if err := WrapTimeFields(&timeValue, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, timeValue.Hour(), localHour)

	times := []time.Time{time.Now().UTC()}
	assert.Equal(t, times[0].Hour(), utcHour)

	if err := WrapTimeFields(times, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, times[0].Hour(), localHour)

	timeStruct := struct {
		T time.Time
	}{
		T: time.Now().UTC(),
	}
	assert.Equal(t, timeStruct.T.Hour(), utcHour)

	if err := WrapTimeFields(timeStruct, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, timeStruct.T.Hour(), utcHour)

	if err := WrapTimeFields(&timeStruct, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, timeStruct.T.Hour(), localHour)

	timeMap := map[string]time.Time{
		"": time.Now().UTC(),
	}
	assert.Equal(t, timeMap[""].Hour(), utcHour)

	if err := WrapTimeFields(timeMap, time.Time.Local); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, timeMap[""].Hour(), localHour)
}
