package reflectx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToNumber(t *testing.T) {
	tests := []struct {
		data any
		typ  string
		want any
	}{
		{data: "123", typ: "int", want: 123},
		{data: "123", typ: "int8", want: int8(123)},
		{data: 123, typ: "int", want: 123},
		{data: 123, typ: "int8", want: int8(123)},
		{data: json.Number("123"), typ: "int64", want: int64(123)},
	}

	for _, test := range tests {
		got, err := ConvertToNumber(test.data, test.typ)
		if err != nil {
			t.Errorf("ConvertToNumber(%v, %s) = %v, want %v", test.data, test.typ, got, test.want)
			return
		}
		require.Equal(t, test.want, got)
	}
}