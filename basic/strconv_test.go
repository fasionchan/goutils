package basic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		str  string
		typ  string
		want any
	}{
		{str: "123", typ: "int", want: 123},
		{str: "123", typ: "int8", want: int8(123)},
		{str: "123", typ: "int16", want: int16(123)},
		{str: "123", typ: "int32", want: int32(123)},
		{str: "123", typ: "int64", want: int64(123)},

		{str: "123", typ: "uint", want: uint(123)},
		{str: "123", typ: "uint8", want: uint8(123)},
		{str: "123", typ: "uint16", want: uint16(123)},
		{str: "123", typ: "uint32", want: uint32(123)},
		{str: "123", typ: "uint64", want: uint64(123)},

		{str: "123", typ: "float32", want: float32(123)},
		{str: "123", typ: "float64", want: float64(123)},

		{str: "123.456", typ: "float32", want: float32(123.456)},
		{str: "123.456", typ: "float64", want: float64(123.456)},

		// {str: "257", typ: "uint8", want: uint8(1)},
	}

	for _, test := range tests {
		got, err := ParseNumber(test.str, test.typ)
		if err != nil {
			t.Errorf("ParseNumber(%s, %s) = %v", test.str, test.typ, err)
			return
		}

		require.Equal(t, test.want, got)
	}
}
