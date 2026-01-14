package model

import "time"

type ZhaoykCou struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey" json:"id"`
	Name       string    `gorm:"column:zyk_name;size:100;not null" json:"name"`
	SchoolYear int64     `gorm:"column:zyk_school_year" json:"school_year"`
	Semester   string    `gorm:"column:zyk_semester;size:1" json:"semester"`
	Hours      int64     `gorm:"column:zyk_hours" json:"hours"`
	Credit     float64   `gorm:"column:zyk_credit;type:numeric(4,2)" json:"credit"`
	ClassID    int64     `gorm:"column:zyk_class_id;not null" json:"class_id"`
	ExamType   string    `gorm:"column:zyk_exam_type;size:10" json:"exam_type"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
