package types

import "testing"

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		value     float64
		precision int
		trimZero  bool
		want      string
	}{
		{1.23456789, 2, true, "1.23"},
		{1.23456789, 2, false, "1.23"},
		{1.23456789, 3, true, "1.235"},
		{1.23456789, 3, false, "1.235"},
		{1.20000000, 2, true, "1.2"},
		{1.20000000, 2, false, "1.20"},
		{1.20000000, 3, true, "1.2"},
		{1.20000000, 3, false, "1.200"},
		{1.00000000, 2, true, "1"},
		{1.00000000, 2, false, "1.00"},
	}

	for _, test := range tests {
		if got := FormatFloat(test.value, test.precision, test.trimZero); got != test.want {
			t.Errorf("FormatFloat(%f, %d, %t) = %s, want %s", test.value, test.precision, test.trimZero, got, test.want)
		}
	}
}
