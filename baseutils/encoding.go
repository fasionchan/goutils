/*
 * Author: fasion
 * Created time: 2023-11-03 10:39:49
 * Last Modified by: fasion
 * Last Modified time: 2023-11-03 11:15:43
 */

package baseutils

import "encoding/json"

type JsonAny = any
type JsonObject = map[string]JsonAny
type JsonMap = JsonObject
type JsonArray = []JsonAny

type JsonRawMessage struct {
	json.RawMessage
}

func (message JsonRawMessage) UnmarshalTo(dst interface{}) error {
	return json.Unmarshal(message.RawMessage, dst)
}

func (message JsonRawMessage) UnmarshalAsString() (string, error) {
	return ParseJsonRawMessage[string](message)
}

func (message JsonRawMessage) UnmarshalAsInt() (int, error) {
	return ParseJsonRawMessage[int](message)
}

func ParseJsonRawMessageAsPointer[Ptr *Data, Data any](message JsonRawMessage) (Ptr, error) {
	var data Data
	if err := message.UnmarshalTo(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func ParseJsonRawMessage[Data any](message JsonRawMessage) (data Data, err error) {
	err = message.UnmarshalTo(&data)
	return
}

type JsonRawMessageMap map[string]JsonRawMessage

func (m JsonRawMessageMap) Lookup(key string) (JsonRawMessage, bool) {
	v, ok := m[key]
	return v, ok
}
