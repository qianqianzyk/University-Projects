package model

import "time"

type ZhaoykTeaching struct {
	ID         int64     `gorm:"column:zyk_id;primaryKey" json:"id"`
	TeacherID  string    `gorm:"column:zyk_teacher_id;size:10;not null" json:"teacher_id"`
	CourseID   int64     `gorm:"column:zyk_course_id;not null" json:"course_id"`
	AssignTime time.Time `gorm:"column:zyk_assign_time" json:"assign_time"`
	CreateTime time.Time `gorm:"column:zyk_create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:zyk_update_time;autoUpdateTime" json:"update_time"`
}
