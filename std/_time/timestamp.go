package _time

import (
	"encoding/json"
	"time"
)

type MilliTimestamp int64

func (ts MilliTimestamp) Time() time.Time {
	return time.UnixMilli(int64(ts))
}

type TimeWithUnixMilliJson time.Time

func (t TimeWithUnixMilliJson) Native() time.Time {
	return time.Time(t)
}

func (t TimeWithUnixMilliJson) Timestamp() int64 {
	return t.Native().UnixMilli()
}

func (t TimeWithUnixMilliJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Native().UnixMilli())
}

func (t *TimeWithUnixMilliJson) UnmarshalJSON(data []byte) error {
	var milli int64
	if err := json.Unmarshal(data, &milli); err != nil {
		return err
	}

	*t = TimeWithUnixMilliJson(time.UnixMilli(milli))
	return nil
}

type TimeWithUnixJson time.Time

func (t TimeWithUnixJson) Native() time.Time {
	return time.Time(t)
}

func (t TimeWithUnixJson) Timestamp() int64 {
	return t.Native().Unix()
}

func (t TimeWithUnixJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Native().Unix())
}

func (t *TimeWithUnixJson) UnmarshalJSON(data []byte) error {
	var sec int64
	if err := json.Unmarshal(data, &sec); err != nil {
		return err
	}

	*t = TimeWithUnixJson(time.Unix(sec, 0))
	return nil
}