/*
 * Author: fasion
 * Created time: 2025-12-07 15:45:43
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:18:25
 */

package entity

import "time"

type Relation struct {
	// 实例ID
	InstanceId string `bson:"InstanceId" json:"InstanceId"`

	// 关系元数据
	MetaName string        `bson:"MetaName" json:"MetaName"`
	Meta     *RelationMeta `bson:"-" json:"Meta"`

	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`           // 创建时间
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"` // 最近修改时间

	Deleted bool `bson:"Deleted" json:"Deleted"` // 是否删除
}
