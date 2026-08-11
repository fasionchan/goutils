package jsonx

import (
	"testing"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/stretchr/testify/require"
)

type JsonDecodeAnyToAnyTestCase struct {
	_testing.TestCaseName
	jsonData any
	expertedErr bool
	expectedResult any
}

func (test *JsonDecodeAnyToAnyTestCase) Run(t *testing.T) {
	result, err := JsonDecodeAnyToAny(test.jsonData)
	if test.expertedErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
		require.Equal(t, test.expectedResult, result)
	}
}

func TestJsonDecodeAnyToAny(t *testing.T) {
	type String string
	type Bytes []byte

	_testing.TypedRunNamedTestCases(t, []*JsonDecodeAnyToAnyTestCase{
		{
			TestCaseName: "string",
			jsonData: `{"name": "John", "age": 30}`,
			expectedResult: map[string]any{
				"name": "John",
				"age": float64(30),
			},
		},
		{
			TestCaseName: "[]byte",
			jsonData: []byte(`{"name": "John", "age": 30}`),
			expectedResult: map[string]any{
				"name": "John",
				"age": float64(30),
			},
		},
		{
			TestCaseName: "Custom String",
			jsonData: String(`{"name": "John", "age": 30}`),
			expectedResult: map[string]any{
				"name": "John",
				"age": float64(30),
			},
		},
		{
			TestCaseName: "Custom Bytes",
			jsonData: Bytes(`{"name": "John", "age": 30}`),
			expectedResult: map[string]any{
				"name": "John",
				"age": float64(30),
			},
		},
		{
			TestCaseName: "Invalid JSON",
			jsonData: "invalid JSON",
			expertedErr: true,
		},
	})
}