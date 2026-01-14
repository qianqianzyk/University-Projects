package model

import "time"

type ZhaoykProvince struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:zyk_name;size:100;unique;not null" json:"name"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
