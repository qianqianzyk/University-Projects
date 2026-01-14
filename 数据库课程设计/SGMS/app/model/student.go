package model

import "time"

type ZhaoykStu struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey" json:"id"`
	StudentID  string    `gorm:"column:zyk_student_id;size:12;unique;not null" json:"student_id"`
	Password   string    `gorm:"column:zyk_password;size:128;not null" json:"-"`
	Name       string    `gorm:"column:zyk_name;size:50;not null" json:"name"`
	Gender     string    `gorm:"column:zyk_gender;size:1" json:"gender,omitempty"`
	Age        int64     `gorm:"column:zyk_age" json:"age"`
	CityID     int64     `gorm:"column:zyk_city_id" json:"city_id"`
	ClassID    int64     `gorm:"column:zyk_class_id;not null" json:"class_id"`
	GPA        float64   `gorm:"column:zyk_gpa;type:numeric(4,2);default:0" json:"gpa"`
	Credits    float64   `gorm:"column:zyk_credits;type:numeric(5,2);default:0" json:"credits"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
