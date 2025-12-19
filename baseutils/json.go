/*
 * Author: fasion
 * Created time: 2025-12-19 15:17:53
 * Last Modified by: fasion
 * Last Modified time: 2025-12-19 15:18:03
 */

package baseutils

import "errors"

type JsonRawStrMessage string

func (msg JsonRawStrMessage) MarshalJSON() ([]byte, error) {
	if msg == "" {
		return []byte("null"), nil
	}
	return []byte(msg), nil
}

func (msg *JsonRawStrMessage) UnmarshalJSON(data []byte) error {
	if msg == nil {
		return errors.New("JsonRawStrMessage is nil")
	}

	*msg = JsonRawStrMessage(data)
	return nil
}
