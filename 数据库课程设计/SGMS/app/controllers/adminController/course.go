package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetCourseReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
	ID      int64 `form:"id" binding:"required"`
}

type GetCourseByAdminResp struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	ClassID    int64   `json:"class_id"`
	SchoolYear int64   `json:"school_year"`
	Semester   string  `json:"semester"`
	Credit     float64 `json:"credit"`
	Hours      int64   `json:"hours"`
	ExamType   string  `json:"exam_type"`
}

func GetCourse(c *gin.Context) {
	var req GetCourseReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	var course model.ZhaoykCou
	if err := database.DB.Raw("SELECT * FROM zhaoyk_cou WHERE zyk_id = ?", req.ID).Scan(&course).Error; err != nil || course.ID == 0 {
		apiException.AbortWithException(c, apiException.CourseNotFoundError, err)
		return
	}

	utils.JsonSuccessResponse(c, GetCourseByAdminResp{
		ID:         course.ID,
		Name:       course.Name,
		ClassID:    course.ClassID,
		SchoolYear: course.SchoolYear,
		Semester:   course.Semester,
		Credit:     course.Credit,
		Hours:      course.Hours,
		ExamType:   course.ExamType,
	})
}

type DeleteCourseReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

func DeleteCourse(c *gin.Context) {
	var req DeleteCourseReq
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	if err := database.DB.Exec("DELETE FROM zhaoyk_cou WHERE zyk_id = ?", req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateCourseReq struct {
	AdminID    int64   `json:"admin_id"`
	Name       string  `json:"name"`
	SchoolYear int64   `json:"school_year"`
	Semester   string  `json:"semester"`
	Hours      int64   `json:"hours"`
	Credit     float64 `json:"credit"`
	ClassID    int64   `json:"class_id"`
	ExamType   string  `json:"exam_type"`
}

func CreateCourse(c *gin.Context) {
	var req CreateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	var newCourseID int64
	insertSQL := `
		INSERT INTO zhaoyk_cou (zyk_name, zyk_school_year, zyk_semester, zyk_hours, zyk_credit, zyk_class_id, zyk_exam_type)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING zyk_id
	`
	if err := database.DB.Raw(insertSQL, req.Name, req.SchoolYear, req.Semester, req.Hours, req.Credit, req.ClassID, req.ExamType).
		Scan(&newCourseID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	var students []model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_class_id = ?", req.ClassID).Scan(&students).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	for _, stu := range students {
		if err := database.DB.Exec(
			"INSERT INTO zhaoyk_score (zyk_student_id, zyk_course_id, zyk_score) VALUES (?, ?, NULL)",
			stu.StudentID, newCourseID,
		).Error; err != nil {
			apiException.AbortWithException(c, apiException.DatabaseError, err)
			return
		}
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateCourseReq struct {
	AdminID    int64   `json:"admin_id"`
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	SchoolYear int64   `json:"school_year"`
	Semester   string  `json:"semester"`
	Hours      int64   `json:"hours"`
	Credit     float64 `json:"credit"`
	ClassID    int64   `json:"class_id"`
	ExamType   string  `json:"exam_type"`
}

func UpdateCourse(c *gin.Context) {
	var req UpdateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	if err := database.DB.Exec("UPDATE zhaoyk_cou SET zyk_name = ?, zyk_school_year = ?, zyk_semester = ?, zyk_credit = ?, zyk_hours = ?, zyk_class_id = ?, zyk_exam_type = ? WHERE zyk_id = ?",
		req.Name, req.SchoolYear, req.Semester, req.Credit, req.Hours, req.ClassID, req.ExamType, req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetCourseScoreDistributionReq struct {
	AdminID    int64 `form:"admin_id"`
	CourseID   int64 `form:"course_id"`
	SchoolYear int64 `form:"school_year"`
}

type ScoreDistributionItem struct {
	ScoreRange string `json:"score_range"`
	Count      int64  `json:"count"`
}

type GetCourseScoreDistributionResp struct {
	Distribution []ScoreDistributionItem `json:"distribution"`
}

func GetCourseScoreDistribution(c *gin.Context) {
	var req GetCourseScoreDistributionReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	type CourseScoreDistributionItem struct {
		ScoreRange string `gorm:"column:zyk_score_range"`
		Count      int64  `gorm:"column:zyk_count"`
	}

	var courseDistributions []CourseScoreDistributionItem
	if err := database.DB.Raw("SELECT * FROM zhaoyk_select_course_score_distribution(?, ?)", req.CourseID, req.SchoolYear).Scan(&courseDistributions).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	resp := GetCourseScoreDistributionResp{
		Distribution: make([]ScoreDistributionItem, len(courseDistributions)),
	}
	for i, distribution := range courseDistributions {
		resp.Distribution[i] = ScoreDistributionItem{
			ScoreRange: distribution.ScoreRange,
			Count:      distribution.Count,
		}
	}

	utils.JsonSuccessResponse(c, resp)
}

type GetCourseListReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type CourseListItem struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	ClassID    int64   `json:"class_id"`
	SchoolYear int64   `json:"school_year"`
	Semester   string  `json:"semester"`
	Credit     float64 `json:"credit"`
	Hours      int64   `json:"hours"`
	ExamType   string  `json:"exam_type"`
}

type GetCoursesResp struct {
	List []CourseListItem `json:"list"`
}

func GetCourseList(c *gin.Context) {
	var req GetCourseListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var admin model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&admin).Error; err != nil || admin.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !admin.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	var courseList []model.ZhaoykCou
	if err := database.DB.Raw("SELECT * FROM zhaoyk_cou").Scan(&courseList).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	courses := make([]CourseListItem, 0, len(courseList))
	for _, course := range courseList {
		courses = append(courses, CourseListItem{
			ID:         course.ID,
			Name:       course.Name,
			ClassID:    course.ClassID,
			SchoolYear: course.SchoolYear,
			Semester:   course.Semester,
			Credit:     course.Credit,
			Hours:      course.Hours,
			ExamType:   course.ExamType,
		})
	}

	utils.JsonSuccessResponse(c, GetCoursesResp{List: courses})
}
