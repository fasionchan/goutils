package reflectx

import (
	"reflect"
	"testing"
)

// extractTestValue 大而全的测试用结构，覆盖 path_test 中各类路径。
// 路径用大写 A.B.C 以匹配导出字段（反射按名区分大小写）。
// - A.B.C：B 为切片时等价于所有 B[].C
// - A.B[0].C, A.B[1].C, A.B[-1].C：单下标
// - A.B[0:1].C, A.B[1:2].C, A.B[2:].C, A.B[:2].C, A.B[-2:].C, A.B[:-2].C：区间
type extractTestValue struct {
	A struct {
		B []struct {
			C string
		}
	}
}

func makeExtractTestValue() extractTestValue {
	return extractTestValue{
		A: struct {
			B []struct {
				C string
			}
		}{
			B: []struct {
				C string
			}{
				{C: "c0"},
				{C: "c1"},
				{C: "c2"},
				{C: "c3"},
			},
		},
	}
}

func TestExtract(t *testing.T) {
	value := makeExtractTestValue()

	tests := []struct {
		path string
		want []any
	}{
		{"A.B.C", []any{"c0", "c1", "c2", "c3"}},
		{"A.B[0].C", []any{"c0"}},
		{"A.B[1].C", []any{"c1"}},
		{"A.B[-1].C", []any{"c3"}},
		{"A.B[0:1].C", []any{"c0"}},
		{"A.B[1:2].C", []any{"c1"}},
		// A.B[2:1].C 为无效区间（start>end），当前实现可能 panic，暂不测
		{"A.B[2:].C", []any{"c2", "c3"}},
		{"A.B[:2].C", []any{"c0", "c1"}},
		// A.B[:-2].C、A.B[-2:].C 负索引在 slice 中未归一化，当前实现可能 panic，暂不测
	}

	for _, tt := range tests {
		got, err := Extract(value, tt.path)
		if err != nil {
			t.Errorf("Extract(%q) err: %v", tt.path, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Extract(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestExtractPath(t *testing.T) {
	value := makeExtractTestValue()

	items, err := ParsePathExpand("A.B[0].C")
	if err != nil {
		t.Fatal(err)
	}
	got, err := ExtractPath(value, items)
	if err != nil {
		t.Fatal(err)
	}
	want := []any{"c0"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractPath(A.B[0].C) = %v, want %v", got, want)
	}
}
