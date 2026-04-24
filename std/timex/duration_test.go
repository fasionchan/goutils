package timex

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseIso8601Duration(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		expected Duration
		wantErr  bool
	}{
		{
			name:     "empty string",
			s:        "",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "positive duration",
			s:        "P1Y2M3DT4H5M6S",
			expected: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			wantErr:  false,
		},
		{
			name:     "negative duration",
			s:        "-P1Y2M3DT4H5M6S",
			expected: -Year - 2*Month - 3*Day - 4*Hour - 5*Minute - 6*Second,
			wantErr:  false,
		},
		{
			name:     "only time part",
			s:        "PT1H30M",
			expected: Hour + 30*Minute,
			wantErr:  false,
		},
		{
			name:     "only date part",
			s:        "P1Y",
			expected: Year,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseIso8601Duration(tc.s)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseIso8601Duration() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("ParseIso8601Duration() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_Duration(t *testing.T) {
	testCases := []struct {
		name     string
		duration Duration
		expected time.Duration
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: 0,
		},
		{
			name:     "positive duration",
			duration: Hour + 30*Minute,
			expected: time.Hour + 30*time.Minute,
		},
		{
			name:     "negative duration",
			duration: -Hour - 30*Minute,
			expected: -time.Hour - 30*time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.duration.Duration()
			if actual != tc.expected {
				t.Errorf("Duration.Duration() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_Parts(t *testing.T) {
	testCases := []struct {
		name     string
		duration Duration
		expected [6]int
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: [6]int{0, 0, 0, 0, 0, 0},
		},
		{
			name:     "positive duration",
			duration: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expected: [6]int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "negative duration",
			duration: -Year - 2*Month - 3*Day - 4*Hour - 5*Minute - 6*Second,
			expected: [6]int{-1, -2, -3, -4, -5, -6},
		},
		{
			name:     "only time part",
			duration: 4*Hour + 5*Minute + 6*Second,
			expected: [6]int{0, 0, 0, 4, 5, 6},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y, m, d, h, min, s := tc.duration.Parts()
			actual := [6]int{y, m, d, h, min, s}
			if actual != tc.expected {
				t.Errorf("Duration.Parts() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_YearMonthDayDuration(t *testing.T) {
	testCases := []struct {
		name        string
		duration    Duration
		expectedY   int
		expectedM   int
		expectedD   int
		expectedDur time.Duration
	}{
		{
			name:        "zero duration",
			duration:    0,
			expectedY:   0,
			expectedM:   0,
			expectedD:   0,
			expectedDur: 0,
		},
		{
			name:        "positive duration",
			duration:    Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expectedY:   1,
			expectedM:   2,
			expectedD:   3,
			expectedDur: 4*time.Hour + 5*time.Minute + 6*time.Second,
		},
		{
			name:        "only time part",
			duration:    4*Hour + 5*Minute + 6*Second,
			expectedY:   0,
			expectedM:   0,
			expectedD:   0,
			expectedDur: 4*time.Hour + 5*time.Minute + 6*time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y, m, d, dur := tc.duration.YearMonthDayDuration()
			if y != tc.expectedY || m != tc.expectedM || d != tc.expectedD || dur != tc.expectedDur {
				t.Errorf("Duration.YearMonthDayDuration() = (%d, %d, %d, %v), expected (%d, %d, %d, %v)",
					y, m, d, dur, tc.expectedY, tc.expectedM, tc.expectedD, tc.expectedDur)
			}
		})
	}
}

func TestDuration_AddTo(t *testing.T) {
	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name     string
		duration Duration
		expected time.Time
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: baseTime,
		},
		{
			name:     "positive duration",
			duration: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expected: time.Date(2021, 3, 4, 4, 5, 6, 0, time.UTC),
		},
		{
			name:     "negative duration",
			duration: -Year - 2*Month - 3*Day - 4*Hour - 5*Minute - 6*Second,
			expected: time.Date(2018, 10, 28, 19, 54, 54, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.duration.AddTo(baseTime)
			if !actual.Equal(tc.expected) {
				t.Errorf("Duration.AddTo() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_SubFrom(t *testing.T) {
	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name     string
		duration Duration
		expected time.Time
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: baseTime,
		},
		{
			name:     "positive duration",
			duration: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expected: time.Date(2018, 10, 28, 19, 54, 54, 0, time.UTC),
		},
		{
			name:     "negative duration",
			duration: -Year - 2*Month - 3*Day - 4*Hour - 5*Minute - 6*Second,
			expected: time.Date(2021, 3, 4, 4, 5, 6, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.duration.SubFrom(baseTime)
			if !actual.Equal(tc.expected) {
				t.Errorf("Duration.SubFrom() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_Iso8601String(t *testing.T) {
	testCases := []struct {
		name     string
		duration Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "P0D",
		},
		{
			name:     "positive duration",
			duration: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expected: "P1Y2M3DT4H5M6S",
		},
		{
			name:     "negative duration",
			duration: -Year - 2*Month - 3*Day - 4*Hour - 5*Minute - 6*Second,
			expected: "-P1Y2M3DT4H5M6S",
		},
		{
			name:     "only time part",
			duration: 4*Hour + 5*Minute + 6*Second,
			expected: "PT4H5M6S",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.duration.Iso8601String()
			if actual != tc.expected {
				t.Errorf("Duration.Iso8601String() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		duration Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: `"P0D"`,
		},
		{
			name:     "positive duration",
			duration: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			expected: `"P1Y2M3DT4H5M6S"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.duration)
			if err != nil {
				t.Errorf("Duration.MarshalJSON() error = %v", err)
				return
			}
			if string(actual) != tc.expected {
				t.Errorf("Duration.MarshalJSON() = %v, expected %v", string(actual), tc.expected)
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		jsonStr  string
		expected Duration
		wantErr  bool
	}{
		{
			name:     "zero duration",
			jsonStr:  `"P0D"`,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "positive duration",
			jsonStr:  `"P1Y2M3DT4H5M6S"`,
			expected: Year + 2*Month + 3*Day + 4*Hour + 5*Minute + 6*Second,
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			jsonStr:  `P1Y2M3DT4H5M6S`,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid duration",
			jsonStr:  `"invalid"`,
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual Duration
			err := json.Unmarshal([]byte(tc.jsonStr), &actual)
			if (err != nil) != tc.wantErr {
				t.Errorf("Duration.UnmarshalJSON() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("Duration.UnmarshalJSON() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}

func TestDuration_RandomBetween(t *testing.T) {
	testCases := []struct {
		name string
		d1   Duration
		d2   Duration
	}{{
		name: "positive range",
		d1:   Hour,
		d2:   2 * Hour,
	},
		{
			name: "negative range",
			d1:   -2 * Hour,
			d2:   -Hour,
		},
		{
			name: "mixed range",
			d1:   -Hour,
			d2:   Hour,
		},
		{
			name: "equal values",
			d1:   Hour,
			d2:   Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min := min(int64(tc.d1), int64(tc.d2))
			max := max(int64(tc.d1), int64(tc.d2))

			// 测试多次以确保结果在范围内
			for i := 0; i < 100; i++ {
				result := tc.d1.RandomBetween(tc.d2)
				resultInt := int64(result)
				if resultInt < min || resultInt > max {
					t.Errorf("Duration.RandomBetween() = %v, expected between %v and %v", result, min, max)
				}
			}
		})
	}
}
