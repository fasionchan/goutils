/*
 * Author: fasion
 * Created time: 2025-12-07 15:41:30
 * Last Modified by: fasion
 * Last Modified time: 2025-12-07 16:15:21
 */

package entity

import "time"

type CommonModel struct {
	CreatedTime      time.Time `bson:"CreatedTime" json:"CreatedTime"`
	LastModifiedTime time.Time `bson:"LastModifiedTime" json:"LastModifiedTime"`

	Deleted bool `bson:"Deleted" json:"Deleted"`
}

type CommonMetaModel struct {
	Name  string `bson:"Name" json:"Name"`
	Title string `bson:"Title" json:"Title"`
	Desc  string `bson:"Desc" json:"Desc"`
}
