/*
 * Author: fasion
 * Created time: 2025-12-07 12:36:50
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:17:50
 */

package entity

import "time"

type Entity struct {
	// 实例ID
	InstanceId string `bson:"InstanceId" json:"InstanceId"`

	// 实体元数据
	MetaName string      `bson:"MetaName" json:"MetaName"`
	Meta     *EntityMeta `bson:"-" json:"Meta"`

	// 属性列表
	Attributes Attributes `bson:"Attributes" json:"Attributes"`

	SearchableText string `bson:"SearchableText" json:"SearchableText"`

	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`           // 创建时间
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"` // 最近修改时间

	Deleted bool `bson:"Deleted" json:"Deleted"` // 是否删除
}
