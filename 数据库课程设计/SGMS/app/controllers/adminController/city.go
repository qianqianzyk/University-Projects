package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetCityListReq struct {
	AdminID int64 `form:"admin_id" binding:"required"` // 管理员ID
}

type CityItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ProvinceID int64  `json:"province_id"`
}

func GetCityList(c *gin.Context) {
	var req GetCityListReq
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

	var cityList []model.ZhaoykCity
	if err := database.DB.Raw("SELECT * FROM zhaoyk_city").Scan(&cityList).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	result := make([]gin.H, len(cityList))
	for i, city := range cityList {
		result[i] = gin.H{
			"id":          city.ID,
			"name":        city.Name,
			"province_id": city.ProvinceID,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"cities": result,
	})
}

type GetCityStudentCountReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type ProvinceCityStudentCount struct {
	ProvinceID   int64  `gorm:"column:zyk_province_id" json:"province_id"`
	ProvinceName string `gorm:"column:zyk_province_name" json:"province_name"`
	CityID       int64  `gorm:"column:zyk_city_id" json:"city_id"`
	CityName     string `gorm:"column:zyk_city_name" json:"city_name"`
	StudentCount int64  `gorm:"column:zyk_student_count" json:"student_count"`
}

func GetCityStudentCount(c *gin.Context) {
	var req GetCityStudentCountReq
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

	var counts []ProvinceCityStudentCount
	if err := database.DB.Raw("SELECT * FROM zhaoyk_province_city_student_count_view").Scan(&counts).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"list": counts,
	})
}
