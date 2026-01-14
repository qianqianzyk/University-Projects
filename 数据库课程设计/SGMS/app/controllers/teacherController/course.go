package teacherController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetTeacherCourseByTeacherReq struct {
	TeacherID int64 `form:"teacher_id" binding:"required"`
	Year      int64 `form:"year"`
	Semester  int64 `form:"semester"`
}

type teacherCourseList struct {
	TeacherID   string  `json:"teacher_id" gorm:"column:zyk_teacher_id"`
	TeacherName string  `json:"teacher_name" gorm:"column:zyk_teacher_name"`
	CourseID    int64   `json:"course_id" gorm:"column:zyk_course_id"`
	CourseName  string  `json:"course_name" gorm:"column:zyk_course_name"`
	ClassID     int64   `json:"class_id" gorm:"column:zyk_class_id"`
	ClassName   string  `json:"class_name" gorm:"column:zyk_class_name"`
	SchoolYear  int64   `json:"school_year" gorm:"column:zyk_school_year"`
	Semester    int64   `json:"semester" gorm:"-"`
	Hours       int64   `json:"hours" gorm:"column:zyk_hours"`
	ExamType    string  `json:"exam_type" gorm:"column:zyk_exam_type"`
	Credit      float64 `json:"credit" gorm:"column:zyk_credit"`
	SemesterStr string  `json:"-" gorm:"column:zyk_semester"`
}

func GetTeacherCourseByTeacher(c *gin.Context) {
	var req GetTeacherCourseByTeacherReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	courses := make([]teacherCourseList, 0)
	err := database.DB.Raw(
		"SELECT * FROM Zhaoyk_select_teacher_courses(?, ?, ?)",
		req.TeacherID, req.Year, string(rune('0'+req.Semester)),
	).Scan(&courses).Error
	if err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	courseList := make([]gin.H, 0)
	for _, course := range courses {
		semester := int64(course.SemesterStr[0] - '0')

		teachings := make([]model.ZhaoykTeaching, 0)
		if err := database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.CourseID).Scan(&teachings).Error; err != nil {
			apiException.AbortWithException(c, apiException.DatabaseError, err)
			return
		}

		teachers := make([]gin.H, 0)
		for _, teaching := range teachings {
			var tea model.ZhaoykTea
			if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", teaching.TeacherID).Scan(&tea).Error; err != nil || tea.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
				return
			}
			teachers = append(teachers, gin.H{
				"teacher_id":   tea.ID,
				"teacher_name": tea.Name,
			})
		}

		courseList = append(courseList, gin.H{
			"course_id":   course.CourseID,
			"course_name": course.CourseName,
			"class_id":    course.ClassID,
			"class_name":  course.ClassName,
			"school_year": course.SchoolYear,
			"semester":    semester,
			"hours":       course.Hours,
			"exam_type":   course.ExamType,
			"credits":     course.Credit,
			"teachers":    teachers,
		})
	}

	utils.JsonSuccessResponse(c, gin.H{
		"courses": courseList,
	})
}

type GetCourseAvgScoreReq struct {
	TeacherID int64 `form:"teacher_id" binding:"required"`
}

type teacherCourseAvgScore struct {
	TeacherID  string  `json:"teacher_id" gorm:"column:zyk_teacher_id"`
	CourseName string  `json:"course_name" gorm:"column:zyk_course_name"`
	ClassName  string  `json:"class_name" gorm:"column:zyk_class_name"`
	SchoolYear int64   `json:"school_year" gorm:"column:zyk_school_year"`
	AvgScore   float64 `json:"avg_score" gorm:"column:zyk_avg_score"`
}

func GetCourseAvgScore(c *gin.Context) {
	var req GetCourseAvgScoreReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	avgScores := make([]teacherCourseAvgScore, 0)
	if err := database.DB.Raw("SELECT * FROM zhaoyk_select_teacher_course_avg_scores(?)", teacher.ID).Scan(&avgScores).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	result := make([]gin.H, len(avgScores))
	for i, item := range avgScores {
		result[i] = gin.H{
			"course_name": item.CourseName,
			"class_name":  item.ClassName,
			"school_year": item.SchoolYear,
			"avg_score":   item.AvgScore,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"avg_scores": result,
	})
}

type GetStudentByTeacherReq struct {
	TeacherID int64  `form:"teacher_id" binding:"required"`
	StudentID string `form:"student_id" binding:"required"`
}

func GetStudentByTeacher(c *gin.Context) {
	var req GetStudentByTeacherReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.StudentID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	var class model.ZhaoykClass
	if err := database.DB.Raw("SELECT * FROM zhaoyk_class WHERE zyk_id = ?", student.ClassID).Scan(&class).Error; err != nil || class.ID == 0 {
		apiException.AbortWithException(c, apiException.ClassNotFoundError, err)
		return
	}

	var city model.ZhaoykCity
	if err := database.DB.Raw("SELECT * FROM zhaoyk_city WHERE zyk_id = ?", student.CityID).Scan(&city).Error; err != nil || city.ID == 0 {
		apiException.AbortWithException(c, apiException.CityNotFoundError, err)
		return
	}

	resp := gin.H{
		"id":            student.ID,
		"student_id":    student.StudentID,
		"name":          student.Name,
		"gender":        student.Gender,
		"age":           student.Age,
		"class_id":      student.ClassID,
		"class_name":    class.Name,
		"gpa":           student.GPA,
		"city_id":       student.CityID,
		"city_name":     city.Name,
		"total_credits": student.Credits,
	}

	utils.JsonSuccessResponse(c, resp)
}

type GetCourseStudentListReq struct {
	TeacherID int64 `form:"teacher_id" binding:"required"`
	CourseID  int64 `form:"course_id" binding:"required"`
}

type courseStudent struct {
	StudentId   int64   `json:"student_id" gorm:"column:zyk_student_id"`
	StudentName string  `json:"student_name" gorm:"column:zyk_student_name"`
	ClassName   string  `json:"class_name" gorm:"column:zyk_class_name"`
	Score       float64 `json:"score" gorm:"column:zyk_score"`
	Rank        int64   `json:"rank" gorm:"column:zyk_rank"`
}

func GetCourseStudentList(c *gin.Context) {
	var req GetCourseStudentListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	var students []courseStudent
	if err := database.DB.Raw(`
		SELECT 
			zyk_student_id, 
			zyk_student_name, 
			zyk_class_name, 
			zyk_score, 
			zyk_rank 
		FROM zhaoyk_select_students_by_teacher_course(?, ?)`,
		teacher.TeacherID, req.CourseID,
	).Scan(&students).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	result := make([]gin.H, len(students))
	for i, stu := range students {
		result[i] = gin.H{
			"student_id":   stu.StudentId,
			"student_name": stu.StudentName,
			"class_name":   stu.ClassName,
			"score":        stu.Score,
			"rank":         stu.Rank,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"students": result,
	})
}

type SetStudentScoreReq struct {
	TeacherID int64   `json:"teacher_id"`
	StudentID int64   `json:"student_id"`
	CourseID  int64   `json:"course_id"`
	Score     float64 `json:"score"`
}

func SetStudentScore(c *gin.Context) {
	var req SetStudentScoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.StudentID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	var course model.ZhaoykCou
	if err := database.DB.Raw("SELECT * FROM zhaoyk_cou WHERE zyk_id = ?", req.CourseID).Scan(&course).Error; err != nil || course.ID == 0 {
		apiException.AbortWithException(c, apiException.CourseNotFoundError, err)
		return
	}

	var score model.ZhaoykScore
	if err := database.DB.Raw(
		"UPDATE zhaoyk_score SET zyk_score = ? WHERE zyk_student_id = ? AND zyk_course_id = ? RETURNING *",
		req.Score, student.StudentID, req.CourseID,
	).Scan(&score).Error; err != nil {
		apiException.AbortWithException(c, apiException.UpdateError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
