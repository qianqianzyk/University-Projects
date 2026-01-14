package studentController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type SelectCourseReq struct {
	StudentID int64 `json:"student_id" binding:"required"`
	CourseID  int64 `json:"course_id" binding:"required"`
}

func SelectCourse(c *gin.Context) {
	var req SelectCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.StudentID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	if err := database.DB.Exec("INSERT INTO zhaoyk_score (zyk_student_id, zyk_course_id) VALUES (?, ?)", student.StudentID, req.CourseID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
