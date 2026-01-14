package teacherController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetTeacherReq struct {
	ID int64 `form:"id"`
}

func GetTeacher(c *gin.Context) {
	var req GetTeacherReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var teacher model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.ID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	resp := gin.H{
		"id":         teacher.ID,
		"teacher_id": teacher.TeacherID,
		"name":       teacher.Name,
		"gender":     teacher.Gender,
		"age":        teacher.Age,
		"title":      teacher.Title,
		"phone":      teacher.Phone,
		"is_admin":   teacher.IsAdmin,
	}

	utils.JsonSuccessResponse(c, resp)
}
