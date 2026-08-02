package jsonx

import "encoding/json"

type JsonRawMessages []json.RawMessage

func UnmarshalJsonRawMessages(raw json.RawMessage) (msgs JsonRawMessages, err error) {
	err = json.Unmarshal(raw, &msgs)
	return
}
