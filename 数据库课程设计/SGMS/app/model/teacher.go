package model

import "time"

type ZhaoykTea struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey" json:"id"`
	TeacherID  string    `gorm:"column:zyk_teacher_id;size:10;unique;not null" json:"teacher_id"`
	Password   string    `gorm:"column:zyk_password;size:128;not null" json:"password"`
	Name       string    `gorm:"column:zyk_name;size:50;not null" json:"name"`
	Gender     string    `gorm:"column:zyk_gender;size:1" json:"gender"`
	Age        int       `gorm:"column:zyk_age" json:"age"`
	Title      string    `gorm:"column:zyk_title;size:50" json:"title"`
	Phone      string    `gorm:"column:zyk_phone;size:20" json:"phone"`
	IsAdmin    bool      `gorm:"column:zyk_is_admin;default:false" json:"is_admin"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
