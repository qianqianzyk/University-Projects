package studentController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetClassCourseReq struct {
	ClassID   int64 `form:"class_id" binding:"required"`
	Year      int64 `form:"year" binding:"required"`
	Semester  int64 `form:"semester" binding:"required"`
	StudentID int64 `form:"student_id"`
}

type TeacherList struct {
	TeacherID   int64  `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
}

type CourseList struct {
	CourseID   int64         `json:"course_id"`
	CourseName string        `json:"course_name"`
	ClassID    int64         `json:"class_id"`
	ClassName  string        `json:"class_name"`
	Credits    float64       `json:"credits"`
	SchoolYear int64         `json:"school_year"`
	Semester   int64         `json:"semester"`
	Hours      int64         `json:"hours"`
	ExamType   string        `json:"exam_type"`
	Teachers   []TeacherList `json:"teachers"`
}

type RawCourse struct {
	CourseID    int64   `gorm:"column:zyk_course_id"`
	CourseName  string  `gorm:"column:zyk_course_name"`
	SchoolYear  int64   `gorm:"column:zyk_school_year"`
	Semester    int64   `gorm:"column:zyk_semester"`
	Credit      float64 `gorm:"column:zyk_credit"`
	Hours       int64   `gorm:"column:zyk_hours"`
	ExamType    string  `gorm:"column:zyk_exam_type"`
	TeacherName string  `gorm:"column:zyk_teacher_name"`
	AssignTime  string  `gorm:"column:zyk_assign_time"`
}

func GetClassCourse(c *gin.Context) {
	var req GetClassCourseReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var class model.ZhaoykClass
	if err := database.DB.Raw("SELECT * FROM zhaoyk_class WHERE zyk_id = ?", req.ClassID).Scan(&class).Error; err != nil || class.ID == 0 {
		apiException.AbortWithException(c, apiException.ClassNotFoundError, err)
		return
	}

	rawCourses := make([]RawCourse, 0)
	if err := database.DB.Raw("SELECT * FROM zhaoyk_select_class_course_schedule(?,?,?)",
		req.ClassID, req.Year, string(rune('0'+req.Semester))).Scan(&rawCourses).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	var scores []model.ZhaoykScore
	if req.StudentID != 0 {
		var student model.ZhaoykStu
		if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.StudentID).Scan(&student).Error; err != nil || student.ID == 0 {
			apiException.AbortWithException(c, apiException.UserNotFoundError, err)
			return
		}
		if student.ClassID != req.ClassID {
			apiException.AbortWithException(c, apiException.StudentNotInClassError, nil)
			return
		}
		database.DB.Raw("SELECT * FROM zhaoyk_score WHERE zyk_student_id = ?", student.StudentID).Scan(&scores)
	}

	finalCourseList := make([]CourseList, 0)
	seen := make(map[int64]bool)
	for _, course := range rawCourses {
		if seen[course.CourseID] {
			continue
		}
		seen[course.CourseID] = true
		skip := false
		for _, score := range scores {
			if score.CourseID == course.CourseID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		var teachings []model.ZhaoykTeaching
		database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.CourseID).Scan(&teachings)

		teachers := make([]TeacherList, len(teachings))
		for i, t := range teachings {
			var teacher model.ZhaoykTea
			if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", t.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
				return
			}
			teachers[i] = TeacherList{
				TeacherID:   teacher.ID,
				TeacherName: teacher.Name,
			}
		}

		c := CourseList{
			CourseID:   course.CourseID,
			CourseName: course.CourseName,
			ClassID:    class.ID,
			ClassName:  class.Name,
			Credits:    course.Credit,
			SchoolYear: course.SchoolYear,
			Semester:   course.Semester,
			Hours:      course.Hours,
			ExamType:   course.ExamType,
			Teachers:   teachers,
		}
		finalCourseList = append(finalCourseList, c)
	}

	utils.JsonSuccessResponse(c, gin.H{
		"courses": finalCourseList,
	})
}
