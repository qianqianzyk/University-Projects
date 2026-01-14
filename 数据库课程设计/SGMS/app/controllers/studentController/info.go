package studentController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetStudentInfoReq struct {
	ID string `form:"id"`
}

func GetStudentInfo(c *gin.Context) {
	var req GetStudentInfoReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.ID).Scan(&student).Error; err != nil || student.ID == 0 {
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

type UpdateStudentInfoReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Age    int64  `json:"age,omitempty"`
	CityID int64  `json:"city_id,omitempty"`
}

func UpdateStudentInfo(c *gin.Context) {
	var req UpdateStudentInfoReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.ID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	if req.Name == "" {
		req.Name = student.Name
	}
	if req.Gender == "" {
		req.Gender = student.Gender
	}
	if req.Age == 0 {
		req.Age = student.Age
	}
	if req.CityID == 0 {
		req.CityID = student.CityID
	}

	if err := database.DB.Exec("UPDATE zhaoyk_stu SET zyk_name=?, zyk_gender=?, zyk_age=?, zyk_city_id=? WHERE zyk_id=?",
		req.Name, req.Gender, req.Age, req.CityID, req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.UpdateError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
