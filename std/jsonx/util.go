package jsonx

import (
	"encoding/json"
	"reflect"

	"github.com/fasionchan/goutils/baseutils"
)

var (
	BytesType = reflect.TypeOf([]byte{})
)

func DeepDup[Data any](data Data) (result Data, err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &result)
	return
}

func JsonDecodeAnyToAny(jsonData any) (result any, err error) {
	jsonBytes, ok := jsonData.([]byte)
	if !ok {
		jsonValue := reflect.ValueOf(jsonData)
		if jsonValue.CanConvert(BytesType) {
			jsonBytes = jsonValue.Convert(BytesType).Bytes()
		} else {
			return nil, baseutils.NewBadTypeError("string / []byte", reflect.TypeOf(jsonData).String())
		}
	}

	err = json.Unmarshal(jsonBytes, &result)
	return
}

func JsonEncodeAnyToString(v any) (string, error) {
		data, err := json.Marshal(v)
		return string(data), err
}