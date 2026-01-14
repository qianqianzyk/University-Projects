package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetStudentByAdminReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

type GetStudentByAdminResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int64  `json:"age"`
	StudentID string `json:"student_id"`
	ClassID   int64  `json:"class_id"`
	CityID    int64  `json:"city_id"`
}

func GetStudent(c *gin.Context) {
	var req GetStudentByAdminReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !teacher.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.ID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	utils.JsonSuccessResponse(c, GetStudentByAdminResp{
		ID:        student.ID,
		Name:      student.Name,
		Gender:    student.Gender,
		Age:       student.Age,
		StudentID: student.StudentID,
		ClassID:   student.ClassID,
		CityID:    student.CityID,
	})
}

type DeleteStudentReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

func DeleteStudent(c *gin.Context) {
	var req DeleteStudentReq
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !teacher.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	if err := database.DB.Exec("DELETE FROM zhaoyk_stu WHERE zyk_id = ?", req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateStudentReq struct {
	AdminID   int64  `json:"admin_id"`
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	CityID    int64  `json:"city_id"`
	ClassID   int64  `json:"class_id"`
}

func CreateStudent(c *gin.Context) {
	var req CreateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !teacher.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	password := req.Password
	if password == "" {
		password = ""
	}

	if err := database.DB.Exec(
		"INSERT INTO zhaoyk_stu (zyk_student_id, zyk_password, zyk_name, zyk_gender, zyk_age, zyk_city_id, zyk_class_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.StudentID, password, req.Name, req.Gender, req.Age, req.CityID, req.ClassID,
	).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	var courses []model.ZhaoykCou
	if err := database.DB.Raw("SELECT * FROM zhaoyk_cou WHERE zyk_class_id = ?", req.ClassID).Scan(&courses).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	for _, course := range courses {
		var count int64
		if err := database.DB.Raw(
			"SELECT COUNT(*) FROM zhaoyk_score WHERE zyk_student_id = ? AND zyk_course_id = ?",
			req.StudentID, course.ID,
		).Scan(&count).Error; err != nil {
			apiException.AbortWithException(c, apiException.DatabaseError, err)
			return
		}

		if count == 0 {
			if err := database.DB.Exec(
				"INSERT INTO zhaoyk_score (zyk_student_id, zyk_course_id, zyk_score) VALUES (?, ?, NULL)",
				req.StudentID, course.ID,
			).Error; err != nil {
				apiException.AbortWithException(c, apiException.DatabaseError, err)
				return
			}
		}
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateStudentByAdminReq struct {
	AdminID int64  `json:"admin_id"`
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Age     int    `json:"age"`
	CityID  int64  `json:"city_id"`
	ClassID int64  `json:"class_id"`
}

func UpdateStudent(c *gin.Context) {
	var req UpdateStudentByAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !teacher.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	if err := database.DB.Exec(
		"UPDATE zhaoyk_stu SET zyk_name=?, zyk_gender=?, zyk_age=?, zyk_city_id=?, zyk_class_id=? WHERE zyk_id = ?",
		req.Name, req.Gender, req.Age, req.CityID, req.ClassID, req.ID,
	).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetStudentListReq struct {
	AdminID int64 `form:"admin_id"`
}

type StudentListItem struct {
	ID        int64  `json:"id"`
	StudentID string `json:"student_id"`
	Gender    string `json:"gender"`
	Name      string `json:"name"`
	Age       int64  `json:"age"`
	ClassID   int64  `json:"class_id"`
	CityID    int64  `json:"city_id"`
}

type GetStudentListResp struct {
	List []StudentListItem `json:"list"`
}

func GetStudentList(c *gin.Context) {
	var req GetStudentListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.AdminID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}
	if !teacher.IsAdmin {
		apiException.AbortWithException(c, apiException.UserPermissionError, nil)
		return
	}

	var students []model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu").Scan(&students).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	studentList := make([]StudentListItem, 0, len(students))
	for _, student := range students {
		studentList = append(studentList, StudentListItem{
			ID:        student.ID,
			StudentID: student.StudentID,
			Gender:    student.Gender,
			Name:      student.Name,
			Age:       student.Age,
			ClassID:   student.ClassID,
			CityID:    student.CityID,
		})
	}

	utils.JsonSuccessResponse(c, GetStudentListResp{
		List: studentList,
	})
}
