/*
 * Author: fasion
 * Created time: 2025-12-07 15:27:05
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:24:18
 */

package entity

import "time"

type EntityMeta struct {
	Name  string `bson:"Name" json:"Name"`
	Title string `bson:"Title" json:"Title"`
	Desc  string `bson:"Desc" json:"Desc"`

	Attributes AttributeMetas `bson:"Attributes" json:"Attributes"`

	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`           // 创建时间
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"` // 最近修改时间

	Deleted bool `bson:"Deleted" json:"Deleted"` // 是否删除
}
