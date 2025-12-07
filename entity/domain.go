/*
 * Author: fasion
 * Created time: 2025-12-07 16:42:49
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:43:10
 */

package entity

import "time"

type Domain struct {
	Name  string `bson:"Name" json:"Name"`
	Title string `bson:"Title" json:"Title"`
	Desc  string `bson:"Desc" json:"Desc"`

	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`           // 创建时间
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"` // 最近修改时间

	Deleted bool `bson:"Deleted" json:"Deleted"` // 是否删除
}
