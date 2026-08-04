package jsonx

import "encoding/json"

func DeepDup[Data any](data Data) (result Data, err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &result)
	return
}
