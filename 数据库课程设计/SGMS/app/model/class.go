package model

import "time"

type ZhaoykClass struct {
	ID           int64     `gorm:"column:zyk_id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:zyk_name;size:100;not null" json:"name"`
	DepartmentID int64     `gorm:"column:zyk_department_id;not null" json:"department_id"`
	CreateTime   time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
