/*
 * Author: fasion
 * Created time: 2025-12-07 12:36:44
 * Last Modified by: fasion
 * Last Modified time: 2025-12-08 00:50:46
 */

package entity

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

var (
	IndexPattern = regexp.MustCompile(`\[(\d+)\]`)
)

type EntityAttribute = Attribute
type EntityAttributePtr = AttributePtr
type EntityAttributes = Attributes

type AttributePtr = *Attribute

type Attribute struct {
	Name  string `bson:"Name" json:"Name"`
	Value any    `bson:"Value" json:"Value"`
}

func (attr *Attribute) GetName() string {
	if attr == nil {
		return ""
	}

	return attr.Name
}

func (attr *Attribute) GetValue() any {
	if attr == nil {
		return nil
	}

	return attr.Value
}

func (attr *Attribute) Print() {
	fmt.Println(attr.String())
}

func (attr *Attribute) String() string {
	return fmt.Sprintf("%s: %+v", attr.Name, attr.Value)
}

type Attributes []*Attribute

func NewAttributes(attrs ...*Attribute) Attributes {
	return attrs
}

// FlattenDataAttributes 将嵌套数据转换为扁平的属性列表
// 支持多层嵌套的 slice/array/map/struct
// 格式示例：
//   - name1.name2.name3: value (嵌套 map)
//   - name1.name2[0].name3: value (包含数组)
func FlattenDataAttributes(data any) Attributes {
	return NewAttributes().AppendFromData("", data)
}

func (attrs Attributes) AppendFromData(prefix string, data any) Attributes {
	return attrs.AppendFromValue(prefix, reflect.ValueOf(data))
}

func (attrs Attributes) AppendFromValue(prefix string, value reflect.Value) Attributes {
	// fmt.Println("AppendFromValue", prefix, value.Kind())
	// fmt.Println()

	switch value.Kind() {
	case reflect.Pointer:
		return attrs.AppendFromPointerValue(prefix, value)
	case reflect.Slice, reflect.Array:
		return attrs.AppendFromListValue(prefix, value)
	case reflect.Map:
		return attrs.AppendFromMapValue(prefix, value)
	case reflect.Struct:
		return attrs.AppendFromStructValue(prefix, value)
	case reflect.Interface:
		return attrs.AppendFromValue(prefix, value.Elem())
	default:
		attrs = append(attrs, &Attribute{Name: prefix, Value: value.Interface()})
		return attrs
	}
}

func (attrs Attributes) AppendFromPointerValue(prefix string, value reflect.Value) Attributes {
	if value.IsNil() {
		return attrs
	}
	return attrs.AppendFromValue(prefix, value.Elem())
}

func (attrs Attributes) AppendFromListValue(prefix string, value reflect.Value) Attributes {
	for i := 0; i < value.Len(); i++ {
		attrs = attrs.AppendFromValue(joinPrefixIndex(prefix, i), value.Index(i))
	}
	return attrs
}

func (attrs Attributes) AppendFromMapValue(prefix string, value reflect.Value) Attributes {
	for iter := value.MapRange(); iter.Next(); {
		attrs = attrs.AppendFromValue(joinPrefix(prefix, iter.Key().String()), iter.Value())
	}
	return attrs
}

func (attrs Attributes) AppendFromStructValue(prefix string, value reflect.Value) Attributes {
	// 获取缓存的结构体类型信息
	fieldInfos := getStructAttrFieldInfos(value.Type())

	// 遍历缓存的字段信息
	for _, fieldInfo := range fieldInfos {
		fieldValue := value.Field(fieldInfo.Index)

		// 处理 omitzero 选项
		if fieldInfo.OmitZero && fieldValue.IsZero() {
			continue
		}

		attrs = attrs.AppendFromValue(joinPrefix(prefix, fieldInfo.AttrName), fieldValue)
	}
	return attrs
}

func (attrs Attributes) SortByName() Attributes {
	return stl.Sort(attrs, func(a, b *Attribute) bool {
		return a.GetName() < b.GetName()
	})
}

func (attrs Attributes) Print() {
	stl.ForEach(attrs, AttributePtr.Print)
}

func ParseAttrTag(tag string) (name string, options types.Strings) {
	tag = strings.TrimSpace(tag)
	index := strings.Index(tag, ",")
	if index == -1 {
		return tag, nil
	}
	return tag[:index], types.SplitToStrings(tag[index+1:], ",").TrimSpace().PurgeZero()
}

// joinPrefix 连接前缀和 key，使用 "." 分隔
func joinPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// joinPrefixIndex 连接前缀和数组索引，使用 "[index]" 格式
func joinPrefixIndex(prefix string, index int) string {
	return prefix + "[" + strconv.Itoa(index) + "]"
}

// attrFieldInfo 缓存的属性字段信息
type attrFieldInfo struct {
	Index    int           // 字段索引
	AttrName string        // 属性名称
	Options  types.Strings // 选项
	OmitZero bool          // 是否忽略零值
}

type attrFieldInfos []*attrFieldInfo

// structTypeCache 结构体类型信息缓存
var structTypeCache sync.Map // map[reflect.Type]*structTypeInfo

// getStructTypeInfo 获取结构体类型信息（带缓存）
func getStructAttrFieldInfos(t reflect.Type) attrFieldInfos {
	// 尝试从缓存获取
	if cached, ok := structTypeCache.Load(t); ok {
		return cached.(attrFieldInfos)
	}

	// 解析结构体字段
	info := parseStructType(t)

	// 存入缓存（使用 LoadOrStore 避免并发重复解析）
	actual, _ := structTypeCache.LoadOrStore(t, info)
	return actual.(attrFieldInfos)
}

// parseStructType 解析结构体类型
func parseStructType(t reflect.Type) attrFieldInfos {
	fmt.Println("parseStructType", t.Name(), t.PkgPath(), t.String())
	numField := t.NumField()
	fields := make([]*attrFieldInfo, 0, numField)

	for i := 0; i < numField; i++ {
		field := t.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		attrName, options := ParseAttrTag(field.Tag.Get("attr"))

		// 处理 "-" 跳过标记
		if attrName == "-" {
			continue
		}

		// 默认使用字段名
		if attrName == "" {
			attrName = field.Name
		}

		fields = append(fields, &attrFieldInfo{
			Index:    i,
			AttrName: attrName,
			Options:  options,
			OmitZero: options.Contain("omitzero"),
		})
	}

	return fields
}
