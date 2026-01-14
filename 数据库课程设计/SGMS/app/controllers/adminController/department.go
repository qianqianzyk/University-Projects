package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetDepartmentReq struct {
	AdminID      int64 `form:"admin_id"`
	DepartmentID int64 `form:"department_id"`
}

type GetDepartmentResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func GetDepartment(c *gin.Context) {
	var req GetDepartmentReq
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

	var department model.ZhaoykDepartment
	if err := database.DB.Raw("SELECT * FROM zhaoyk_department WHERE zyk_id = ?", req.DepartmentID).Scan(&department).Error; err != nil || department.ID == 0 {
		apiException.AbortWithException(c, apiException.DepartmentNotFoundError, err)
		return
	}

	utils.JsonSuccessResponse(c, GetDepartmentResp{
		ID:   department.ID,
		Name: department.Name,
	})
}

type DeleteDepartmentReq struct {
	AdminID      int64 `form:"admin_id"`
	DepartmentID int64 `form:"department_id"`
}

func DeleteDepartment(c *gin.Context) {
	var req DeleteDepartmentReq
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

	var department model.ZhaoykDepartment
	if err := database.DB.Raw("DELETE FROM zhaoyk_department WHERE zyk_id = ?", req.DepartmentID).Scan(&department).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateDepartmentReq struct {
	AdminID int64  `json:"admin_id"`
	Name    string `json:"name"`
}

func CreateDepartment(c *gin.Context) {
	var req CreateDepartmentReq
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

	var department model.ZhaoykDepartment
	if err := database.DB.Raw("INSERT INTO zhaoyk_department (zyk_name) VALUES (?)", req.Name).Scan(&department).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateDepartmentReq struct {
	AdminID      int64  `json:"admin_id"`
	DepartmentID int64  `json:"department_id"`
	Name         string `json:"name"`
}

func UpdateDepartment(c *gin.Context) {
	var req UpdateDepartmentReq
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

	var department model.ZhaoykDepartment
	if err := database.DB.Raw("UPDATE zhaoyk_department SET zyk_name = ? WHERE zyk_id = ?", req.Name, req.DepartmentID).Scan(&department).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetDepartmentListReq struct {
	AdminID int64 `form:"admin_id"`
}

type DepartmentListItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetDepartmentListResp struct {
	List []DepartmentListItem `json:"list"`
}

func GetDepartmentList(c *gin.Context) {
	var req GetDepartmentListReq
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

	var departments []model.ZhaoykDepartment
	if err := database.DB.Raw("SELECT * FROM zhaoyk_department").Scan(&departments).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]DepartmentListItem, 0, len(departments))
	for _, d := range departments {
		list = append(list, DepartmentListItem{
			ID:   d.ID,
			Name: d.Name,
		})
	}

	utils.JsonSuccessResponse(c, GetDepartmentListResp{
		List: list,
	})
}
