package reflectx

import (
	"reflect"
	"testing"
)

func Addr[T any](data T) *T {
	return &data
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		path string
		want PathItems
		wantErr bool
	}{
		{
			path: "a.b.c",
			want: PathItems{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
		},
		{
			path: "a.b[0].c",
			want: PathItems{
				{Name: "a"},
				{Name: "b", Indexes: DataIndexes{{index: "0"}}},
				{Name: "c"},
			},
		},
		{
			path: "a.b[1].c",
			want: PathItems{
				{Name: "a"},
				{Name: "b", Indexes: DataIndexes{{index: "1", i: 1}}},
				{Name: "c"},
			},
		},
		{path: "a.b[-1].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "-1", i: -1}}},
			{Name: "c"},
		}},
		{path: "a.b[0:1].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "0:1", ranges: &[2]*int{Addr(0), Addr(1)}}}},
			{Name: "c"},
		}},
		{path: "a.b[1:2].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "1:2", ranges: &[2]*int{Addr(1), Addr(2)}}}},
			{Name: "c"},
		}},
		{path: "a.b[2:1].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "2:1", ranges: &[2]*int{Addr(2), Addr(1)}}}},
			{Name: "c"},
		}},
		{path: "a.b[2:].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "2:", ranges: &[2]*int{Addr(2)}}}},
			{Name: "c"},
		}},
		{path: "a.b[:2].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: ":2", ranges: &[2]*int{nil, Addr(2)}}}},
			{Name: "c"},
		}},
		{path: "a.b[:-2].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: ":-2", ranges: &[2]*int{nil, Addr(-2)}}}},
			{Name: "c"},
		}},
		{path: "a.b[-2:].c", want: PathItems{
			{Name: "a"},
			{Name: "b", Indexes: DataIndexes{{index: "-2:", ranges: &[2]*int{Addr(-2)}}}},
			{Name: "c"},
		}},
		{path: "a.b[-2-1].c", wantErr: true},
	}
	for _, test := range tests {
		got, err := ParsePath(test.path)
		if test.wantErr {
			if err == nil {
				t.Fatalf("ParsePath(%q) expected error, got nil", test.path)
			}
			continue
		}

		if err != nil {
			t.Fatalf("ParsePath(%q) error: %v", test.path, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("ParsePath(%q) = %v, want %v", test.path, got, test.want)
		}
	}
}