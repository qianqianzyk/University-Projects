package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetProvinceListReq struct {
	AdminID int64 `form:"admin_id"`
}

type ProvinceItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetProvinceListResp struct {
	Provinces []ProvinceItem `json:"provinces"`
}

func GetProvinceList(c *gin.Context) {
	var req GetProvinceListReq
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

	var provinces []model.ZhaoykProvince
	if err := database.DB.Raw("SELECT * FROM zhaoyk_province").Scan(&provinces).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]ProvinceItem, 0, len(provinces))
	for _, p := range provinces {
		list = append(list, ProvinceItem{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	utils.JsonSuccessResponse(c, GetProvinceListResp{
		Provinces: list,
	})
}
