package studentController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

type GetTeacherCourseReq struct {
	TeacherID int64 `form:"teacher_id" binding:"required"`
	Year      int64 `form:"year"`
	Semester  int64 `form:"semester"`
}

func GetTeacherCourse(c *gin.Context) {
	var req GetTeacherCourseReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var courses []struct {
		TeacherID   string  `gorm:"column:zyk_teacher_id" json:"teacher_id"`
		TeacherName string  `gorm:"column:zyk_teacher_name" json:"teacher_name"`
		CourseID    int64   `gorm:"column:zyk_course_id" json:"course_id"`
		CourseName  string  `gorm:"column:zyk_course_name" json:"course_name"`
		ClassID     int64   `gorm:"column:zyk_class_id" json:"class_id"`
		ClassName   string  `gorm:"column:zyk_class_name" json:"class_name"`
		SchoolYear  int64   `gorm:"column:zyk_school_year" json:"school_year"`
		Credit      float64 `gorm:"column:zyk_credit" json:"credit"`
		Hours       int64   `gorm:"column:zyk_hours" json:"hours"`
		ExamType    string  `gorm:"column:zyk_exam_type" json:"exam_type"`
		SemesterStr string  `gorm:"column:zyk_semester" json:"-"`
		Semester    int     `json:"semester"`
	}

	semesterStr := strconv.Itoa(int(req.Semester))
	if err := database.DB.Raw("SELECT * FROM Zhaoyk_select_teacher_courses(?, ?, ?)", req.TeacherID, req.Year, semesterStr).Scan(&courses).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	var result []CourseList
	for _, course := range courses {
		var teachings []model.ZhaoykTeaching
		database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.CourseID).Scan(&teachings)

		teachers := make([]TeacherList, 0)
		for _, t := range teachings {
			var teacher model.ZhaoykTea
			database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", t.TeacherID).Scan(&teacher)
			if teacher.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, nil)
				return
			}
			teachers = append(teachers, TeacherList{
				TeacherID:   teacher.ID,
				TeacherName: teacher.Name,
			})
		}

		semInt, err := strconv.Atoi(course.SemesterStr)
		if err != nil {
			semInt = 0
		}

		result = append(result, CourseList{
			CourseID:   course.CourseID,
			CourseName: course.CourseName,
			ClassID:    course.ClassID,
			ClassName:  course.ClassName,
			Credits:    course.Credit,
			SchoolYear: course.SchoolYear,
			Semester:   int64(semInt),
			Hours:      course.Hours,
			ExamType:   course.ExamType,
			Teachers:   teachers,
		})
	}

	utils.JsonSuccessResponse(c, gin.H{
		"courses": result,
	})
}
