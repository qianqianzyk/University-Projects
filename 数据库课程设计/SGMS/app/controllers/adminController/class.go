package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetClassReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
	ClassID int64 `form:"class_id" binding:"required"`
}

type GetClassResp struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DepartmentID int64  `json:"department_id"`
}

func GetClass(c *gin.Context) {
	var req GetClassReq
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

	var class model.ZhaoykClass
	if err := database.DB.Raw("SELECT * FROM zhaoyk_class WHERE zyk_id = ?", req.ClassID).Scan(&class).Error; err != nil || class.ID == 0 {
		apiException.AbortWithException(c, apiException.ClassNotFoundError, err)
		return
	}

	utils.JsonSuccessResponse(c, GetClassResp{
		ID:           class.ID,
		Name:         class.Name,
		DepartmentID: class.DepartmentID,
	})
}

type DeleteClassReq struct {
	AdminID int64 `form:"admin_id"`
	ClassID int64 `form:"class_id"`
}

func DeleteClass(c *gin.Context) {
	var req DeleteClassReq
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

	if err := database.DB.Exec("DELETE FROM zhaoyk_class WHERE zyk_id = ?", req.ClassID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateClassReq struct {
	AdminID      int64  `json:"admin_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	DepartmentID int64  `json:"department_id" binding:"required"`
}

func CreateClass(c *gin.Context) {
	var req CreateClassReq
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

	if err := database.DB.Exec("INSERT INTO zhaoyk_class (zyk_name, zyk_department_id) VALUES (?, ?)", req.Name, req.DepartmentID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateClassReq struct {
	AdminID      int64  `json:"admin_id" binding:"required"`
	ClassID      int64  `json:"class_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	DepartmentID int64  `json:"department_id" binding:"required"`
}

func UpdateClass(c *gin.Context) {
	var req UpdateClassReq
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

	if err := database.DB.Exec("UPDATE zhaoyk_class SET zyk_name = ?, zyk_department_id = ? WHERE zyk_id = ?", req.Name, req.DepartmentID, req.ClassID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetClassCourseListReq struct {
	AdminID    int64 `form:"admin_id"`
	ClassID    int64 `form:"class_id"`
	SchoolYear int64 `form:"school_year"`
	Semester   int64 `form:"semester"`
}

type GetCourseListResp struct {
	List []GetCourseListItem `json:"list"`
}

type GetCourseListItem struct {
	CourseID   int64         `json:"course_id"`
	CourseName string        `json:"course_name"`
	ClassID    int64         `json:"class_id"`
	ClassName  string        `json:"class_name"`
	Credits    float64       `json:"credits"`
	SchoolYear int64         `json:"school_year"`
	Semester   int64         `json:"semester"`
	Hours      int64         `json:"hours"`
	ExamType   string        `json:"exam_type"`
	Teachers   []TeacherItem `json:"teachers"`
}

type TeacherItem struct {
	TeacherID   int64  `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
}

type CourseList struct {
	CourseID    int64   `gorm:"column:zyk_course_id"`
	CourseName  string  `gorm:"column:zyk_course_name"`
	SchoolYear  int64   `gorm:"column:zyk_school_year"`
	Semester    string  `gorm:"column:zyk_semester"`
	Credit      float64 `gorm:"column:zyk_credit"`
	Hours       int64   `gorm:"column:zyk_hours"`
	ExamType    string  `gorm:"column:zyk_exam_type"`
	TeacherName string  `gorm:"column:zyk_teacher_name"`
}

func GetClassCourseList(c *gin.Context) {
	var req GetClassCourseListReq
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

	var class model.ZhaoykClass
	if err := database.DB.Raw("SELECT * FROM zhaoyk_class WHERE zyk_id = ?", req.ClassID).Scan(&class).Error; err != nil || class.ID == 0 {
		apiException.AbortWithException(c, apiException.ClassNotFoundError, err)
		return
	}

	courses := make([]CourseList, 0)
	err := database.DB.Raw("SELECT * FROM zhaoyk_select_class_course_schedule(?,?,?)",
		req.ClassID, req.SchoolYear, string(rune('0'+req.Semester))).Scan(&courses).Error
	if err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	courseList := make([]GetCourseListItem, 0)
	for _, course := range courses {
		teachings := make([]model.ZhaoykTeaching, 0)
		err := database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.CourseID).Scan(&teachings).Error
		if err != nil {
			apiException.AbortWithException(c, apiException.DatabaseError, err)
			return
		}

		teachers := make([]TeacherItem, len(teachings))
		for j, teaching := range teachings {
			var teacher model.ZhaoykTea
			err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", teaching.TeacherID).Scan(&teacher).Error
			if err != nil || teacher.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
				return
			}
			teachers[j] = TeacherItem{
				TeacherID:   teacher.ID,
				TeacherName: teacher.Name,
			}
		}

		courseList = append(courseList, GetCourseListItem{
			CourseID:   course.CourseID,
			CourseName: course.CourseName,
			ClassID:    class.ID,
			ClassName:  class.Name,
			Credits:    course.Credit,
			SchoolYear: course.SchoolYear,
			Semester:   int64(course.Semester[0] - '0'),
			Hours:      course.Hours,
			ExamType:   course.ExamType,
			Teachers:   teachers,
		})
	}

	utils.JsonSuccessResponse(c, GetCourseListResp{
		List: courseList,
	})
}

type GetClassListReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type GetClassListResp struct {
	List []ClassListItem `json:"list"`
}

type ClassListItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DepartmentID int64  `json:"department_id"`
	StudentCount int64  `json:"student_count"`
}

func GetClassList(c *gin.Context) {
	var req GetClassListReq
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

	type ClassWithCount struct {
		ID           int64
		Name         string
		DepartmentID int64
		StudentCount int64
	}

	var classList []ClassWithCount
	if err := database.DB.Raw(`
    	SELECT 
        	c.zyk_id AS id,
        	c.zyk_name AS name,
        	c.zyk_department_id AS department_id,
        	COUNT(s.zyk_id) AS student_count
    	FROM zhaoyk_class c
    	LEFT JOIN zhaoyk_stu s ON c.zyk_id = s.zyk_class_id
    	GROUP BY c.zyk_id, c.zyk_name, c.zyk_department_id
	`).Scan(&classList).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	classes := make([]ClassListItem, 0, len(classList))
	for _, class := range classList {
		classes = append(classes, ClassListItem{
			ID:           class.ID,
			Name:         class.Name,
			DepartmentID: class.DepartmentID,
			StudentCount: class.StudentCount,
		})
	}

	utils.JsonSuccessResponse(c, GetClassListResp{
		List: classes,
	})
}
