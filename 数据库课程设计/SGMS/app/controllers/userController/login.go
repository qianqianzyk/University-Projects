package userController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int64  `json:"user_type"`
}

type LoginResp struct {
	Id       int64 `json:"id"`
	UserType int64 `json:"user_type"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	resp := LoginResp{}
	switch req.UserType {
	case 1:
		var student model.ZhaoykStu
		if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_student_id = ?", req.Username).Scan(&student).Error; err != nil || student == (model.ZhaoykStu{}) {
			apiException.AbortWithException(c, apiException.UserNotFoundError, err)
			return
		}
		if student.Password != req.Password {
			apiException.AbortWithException(c, apiException.UsernameOrPasswordError, nil)
			return
		}
		resp = LoginResp{
			UserType: 1,
			Id:       student.ID,
		}
	case 2:
		var teacher model.ZhaoykTea
		if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", req.Username).Scan(&teacher).Error; err != nil || teacher == (model.ZhaoykTea{}) {
			apiException.AbortWithException(c, apiException.UserNotFoundError, err)
			return
		}
		if teacher.Password != req.Password {
			apiException.AbortWithException(c, apiException.UsernameOrPasswordError, nil)
			return
		}
		resp = LoginResp{
			UserType: 2,
			Id:       teacher.ID,
		}
	case 3:
		var teacher model.ZhaoykTea
		if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", req.Username).Scan(&teacher).Error; err != nil || teacher == (model.ZhaoykTea{}) {
			apiException.AbortWithException(c, apiException.UserNotFoundError, err)
			return
		}
		if !teacher.IsAdmin {
			apiException.AbortWithException(c, apiException.UserPermissionError, nil)
			return
		}
		if teacher.Password != req.Password {
			apiException.AbortWithException(c, apiException.UsernameOrPasswordError, nil)
			return
		}
		resp = LoginResp{
			UserType: 3,
			Id:       teacher.ID,
		}
	default:
		apiException.AbortWithException(c, apiException.UserRoleError, nil)
	}

	utils.JsonSuccessResponse(c, resp)
}
