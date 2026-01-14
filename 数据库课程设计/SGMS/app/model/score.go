package model

import "time"

type ZhaoykScore struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey" json:"id"`
	StudentID  string    `gorm:"column:zyk_student_id;size:12;not null" json:"student_id"`
	CourseID   int64     `gorm:"column:zyk_course_id;not null" json:"course_id"`
	Score      float64   `gorm:"column:zyk_score;type:numeric(5,2)" json:"score"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
