/*
 * Author: fasion
 * Created time: 2025-12-07 15:46:09
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:40:48
 */

package entity

import "time"

type RelationMeta struct {
	Name  string `bson:"Name" json:"Name"`
	Title string `bson:"Title" json:"Title"`
	Desc  string `bson:"Desc" json:"Desc"`

	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`           // 创建时间
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"` // 最近修改时间

	Deleted bool `bson:"Deleted" json:"Deleted"` // 是否删除
}

type RelationMetas []*RelationMeta
