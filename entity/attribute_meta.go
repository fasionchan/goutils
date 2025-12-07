/*
 * Author: fasion
 * Created time: 2025-12-07 15:38:23
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:14:36
 */

package entity

type AttributeMetaPtr = *AttributeMeta

type AttributeMeta struct {
	CommonMetaModel `bson:",inline" json:",inline"`

	// 参数类型
	Type string `bson:"Type" json:"Type"`

	Required   bool `bson:"Required" json:"Required"`     // 是否必填
	Searchable bool `bson:"Searchable" json:"Searchable"` // 是否可搜索
	Vectorized bool `bson:"Vectorized" json:"Vectorized"` // 是否已向量化
}

type AttributeMetas []*AttributeMeta
