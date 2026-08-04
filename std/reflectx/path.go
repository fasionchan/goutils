package reflectx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

// DataIndex 表示路径中的下标，由字符串解析（如 "0"、"-1"、[*]、0:3）。
// - index 为 "*" 或 "" 表示所有元素；"0"、"1" 等为具体下标；"-1" 表示最后一个；0:3 表示 [0, 3) 区间（冒号分隔）。
type DataIndex struct {
	index string

	// 务必将下标解析并保存下来，避免重复解析。
	i int
	ranges *[2]*int
}

func ParseDataIndex(index string) (*DataIndex, error) {
	// 去掉 [ 和 ]
	if strings.HasPrefix(index, "[") {
		if !strings.HasSuffix(index, "]") {
			return nil, fmt.Errorf("invalid data index: %s", index)
		}

		index = index[1 : len(index)-1]
	}

	// 全部元素
	if index == "*" || index == "" {
		return &DataIndex{index: index}, nil
	}

	// 区间 0:3（冒号分隔）
	if strings.Contains(index, ":") {
		parts := strings.SplitN(index, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid data index: %s", index)
		}

		var start, end *int
		if parts[0] != "" {
			i, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid data index %q: %w", index, err)
			}

			start = &i
		}

		if parts[1] != "" {
			i, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid data index %q: %w", index, err)
			}

			end = &i
		}

		return &DataIndex{index: index, ranges: &[2]*int{start, end}}, nil
	}

	// 单个元素（含 "0"、"-1" 等）
	i, err := strconv.Atoi(index)
	if err != nil {
		return nil, fmt.Errorf("invalid data index %q: %w", index, err)
	}

	return &DataIndex{index: index, i: i}, nil
}

func (di *DataIndex) String() string {
	if di == nil {
		return ""
	}
	return "[" + di.index + "]"
}

func (di *DataIndex) StringLite() string {
	if di == nil {
		return ""
	}
	return di.index
}

func (di *DataIndex) IsAll() bool {
	return di != nil && (di.index == "*" || di.index == "")
}

type DataIndexes []*DataIndex

func ParseDataIndexes(indexes string) (DataIndexes, error) {
	// 去掉 [ 和 ]
	if strings.HasPrefix(indexes, "[") {
		if !strings.HasSuffix(indexes, "]") {
			return nil, fmt.Errorf("invalid data indexes: %s", indexes)
		}

		indexes = indexes[1 : len(indexes)-1]
	}

	return stl.MapWithErrorSimplified(strings.Split(indexes, ","), ParseDataIndex)
}

const (
	PathItemTypeAttr = "attr"
	PathItemTypeIndex = "index"
	PathItemTypeIndexes = "indexes"
	PathItemTypeMixed = "mixed"
)

type PathItemPtr = *PathItem

type PathItem struct {
	// - attr: 属性名（如 "a"）
	// - index: 单一下标（如 "[0]"）
	// - indexes: 多个下标（如 "[0,1,2]"）
	// - mixed: 混合类型（如 "attr[0].attr[1]"）
	Type string
	Name  string
	Index *DataIndex
	Indexes DataIndexes
}

func (item *PathItem) IsAtomic() bool {
	if item == nil {
		return false
	}

	return item.Type == PathItemTypeAttr || item.Type == PathItemTypeIndex
}

func (item *PathItem) HasIndex() bool {
	return item.Index != nil
}

func (item *PathItem) String() string {
	if item == nil {
		return ""
	}
	if item.Index == nil {
		return item.Name
	}
	if item.Index.IsAll() {
		return item.Name + "[*]"
	}
	return item.Name + item.Index.String()
}

type PathItems []*PathItem

func (items PathItems) Empty() bool {
	return len(items) == 0
}

func (items PathItems) Strings() types.Strings {
	return  stl.Map(items, PathItemPtr.String)
}

func (items PathItems) String() string {
	return strings.Join(items.Strings(), ".")
}

func (items PathItems) Expand() PathItems {
	result := make(PathItems, 0, len(items))
	for _, item := range items {
		switch item.Type {
		case PathItemTypeAttr:
			result = append(result, item)
		case PathItemTypeIndex:
			result = append(result, item)
		case PathItemTypeIndexes:
			for _, index := range item.Indexes {
				result = append(result, &PathItem{Type: PathItemTypeIndex, Index: index})
			}
		case PathItemTypeMixed, "":
			result = append(result, &PathItem{Type: PathItemTypeAttr, Name: item.Name})
			for _, index := range item.Indexes {
				result = append(result, &PathItem{Type: PathItemTypeIndex, Index: index})
			}
		}
	}
	return result
}

// ParsePath 将路径字符串解析为 PathItems。
// 格式示例：name1.name2、name1.name2[0]、name2[-1]（最后一个）、name2[*] 或 name2[]（全部）。
func ParsePath(s string) (PathItems, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var items PathItems
	for _, part := range strings.Split(s, ".") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		item, err := parsePathItem(part)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func ParsePathExpand(s string) (PathItems, error) {
	items, err := ParsePath(s)
	if err != nil {
		return nil, err
	}
	return items.Expand(), nil
}

func parsePathItem(s string) (*PathItem, error) {
	lb := strings.IndexByte(s, '[')
	if lb < 0 {
		return &PathItem{Name: s}, nil
	}

	name := strings.TrimSpace(s[:lb])
	if name == "" {
		return nil, fmt.Errorf("invalid path item, name is empty: %s", s)
	}

	indexes, err := ParseDataIndexes(s[lb:])
	if err != nil {
		return nil, err
	}

	return &PathItem{Name: name, Indexes: indexes}, nil
}