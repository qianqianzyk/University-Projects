package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetTeacherByAdminReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

type GetTeacherByAdminResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	TeacherID string `json:"teacher_id"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	Title     string `json:"title"`
	Phone     string `json:"phone"`
	IsAdmin   bool   `json:"is_admin"`
}

func GetTeacher(c *gin.Context) {
	var req GetTeacherByAdminReq
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

	var target model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_id = ?", req.ID).Scan(&target).Error; err != nil || target.ID == 0 {
		apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
		return
	}

	utils.JsonSuccessResponse(c, GetTeacherByAdminResp{
		ID:        target.ID,
		Name:      target.Name,
		TeacherID: target.TeacherID,
		Gender:    target.Gender,
		Age:       target.Age,
		Title:     target.Title,
		Phone:     target.Phone,
		IsAdmin:   target.IsAdmin,
	})
}

type DeleteTeacherReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

func DeleteTeacher(c *gin.Context) {
	var req DeleteTeacherReq
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

	var deleteTeacher model.ZhaoykTea
	if err := database.DB.Raw("DELETE FROM zhaoyk_tea WHERE zyk_id = ?", req.ID).Scan(&deleteTeacher).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateTeacherReq struct {
	AdminID   int64  `json:"admin_id"`
	TeacherID string `json:"teacher_id"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	Title     string `json:"title"`
	Phone     string `json:"phone"`
}

func CreateTeacher(c *gin.Context) {
	var req CreateTeacherReq
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

	var teacher model.ZhaoykTea
	if err := database.DB.Raw(
		"INSERT INTO zhaoyk_tea (zyk_teacher_id, zyk_password, zyk_name, zyk_gender, zyk_age, zyk_title, zyk_phone) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.TeacherID, req.Password, req.Name, req.Gender, req.Age, req.Title, req.Phone,
	).Scan(&teacher).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateTeacherReq struct {
	AdminID int64  `json:"admin_id"`
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Age     int    `json:"age"`
	Title   string `json:"title"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
}

func UpdateTeacher(c *gin.Context) {
	var req UpdateTeacherReq
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

	var targetTeacher model.ZhaoykTea
	if err := database.DB.Raw(
		"UPDATE zhaoyk_tea SET zyk_name = ?, zyk_gender = ?, zyk_age = ?, zyk_title = ?, zyk_phone = ?, zyk_is_admin = ? WHERE zyk_id = ?",
		req.Name, req.Gender, req.Age, req.Title, req.Phone, req.IsAdmin, req.ID,
	).Scan(&targetTeacher).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetTeacherCourseListReq struct {
	AdminID    int64 `form:"admin_id"`
	TeacherID  int64 `form:"teacher_id"`
	SchoolYear int64 `form:"school_year"`
	Semester   int64 `form:"semester"`
}

type CourseLists struct {
	CourseID    int64   `gorm:"column:zyk_course_id"`
	CourseName  string  `gorm:"column:zyk_course_name"`
	SchoolYear  int64   `gorm:"column:zyk_school_year"`
	Semester    string  `gorm:"column:zyk_semester"`
	Credit      float64 `gorm:"column:zyk_credit"`
	Hours       int64   `gorm:"column:zyk_hours"`
	ExamType    string  `gorm:"column:zyk_exam_type"`
	TeacherName string  `gorm:"column:zyk_teacher_name"`
	ClassID     int64   `gorm:"column:zyk_class_id"`
	ClassName   string  `gorm:"column:zyk_class_name"`
}

func GetTeacherCourseList(c *gin.Context) {
	var req GetTeacherCourseListReq
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

	var courses []CourseLists
	if err := database.DB.Raw("SELECT * FROM zhaoyk_select_teacher_courses(?, ?, ?)",
		req.TeacherID, req.SchoolYear, string(rune('0'+req.Semester))).Scan(&courses).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	courseList := make([]GetCourseListItem, 0)
	for _, course := range courses {
		var teachings []model.ZhaoykTeaching
		if err := database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.CourseID).Scan(&teachings).Error; err != nil {
			apiException.AbortWithException(c, apiException.DatabaseError, err)
			return
		}
		teachers := make([]TeacherItem, len(teachings))
		for j, teaching := range teachings {
			var tea model.ZhaoykTea
			if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", teaching.TeacherID).Scan(&tea).Error; err != nil || tea.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
				return
			}
			teachers[j] = TeacherItem{
				TeacherID:   tea.ID,
				TeacherName: tea.Name,
			}
		}
		courseList = append(courseList, GetCourseListItem{
			CourseID:   course.CourseID,
			CourseName: course.CourseName,
			ClassID:    course.ClassID,
			ClassName:  course.ClassName,
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

type TeacherCourseView struct {
	TeacherID    string  `gorm:"column:zyk_teacher_id"`
	TeacherName  string  `gorm:"column:zyk_teacher_name"`
	CourseID     int64   `gorm:"column:zyk_course_id"`
	CourseName   string  `gorm:"column:zyk_course_name"`
	SchoolYear   int64   `gorm:"column:zyk_school_year"`
	Semester     string  `gorm:"column:zyk_semester"`
	StudentCount int64   `gorm:"column:zyk_student_count"`
	AvgScore     float64 `gorm:"column:zyk_avg_score"`
}

type GetTeacherStatisticsItem struct {
	TeacherID    string  `json:"teacher_id"`
	TeacherName  string  `json:"teacher_name"`
	CourseID     int64   `json:"course_id"`
	CourseName   string  `json:"course_name"`
	SchoolYear   int64   `json:"school_year"`
	Semester     string  `json:"semester"`
	StudentCount int64   `json:"student_count"`
	AvgScore     float64 `json:"avg_score"`
}

type GetTeacherStatisticsResp struct {
	List []GetTeacherStatisticsItem `json:"list"`
}

type GetTeacherStatisticsReq struct {
	AdminID int64 `form:"admin_id"`
}

func GetTeacherStatistics(c *gin.Context) {
	var req GetTeacherStatisticsReq
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

	var courseViews []TeacherCourseView
	if err := database.DB.Raw("SELECT * FROM zhaoyk_teacher_course_view").Scan(&courseViews).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	courseViewList := make([]GetTeacherStatisticsItem, len(courseViews))
	for i, view := range courseViews {
		courseViewList[i] = GetTeacherStatisticsItem{
			TeacherID:    view.TeacherID,
			TeacherName:  view.TeacherName,
			CourseID:     view.CourseID,
			CourseName:   view.CourseName,
			SchoolYear:   view.SchoolYear,
			Semester:     view.Semester,
			StudentCount: view.StudentCount,
			AvgScore:     view.AvgScore,
		}
	}

	utils.JsonSuccessResponse(c, GetTeacherStatisticsResp{
		List: courseViewList,
	})
}

type GetTeachingReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

type GetTeachingResp struct {
	ID        int64  `json:"id"`
	CourseID  int64  `json:"course_id"`
	TeacherID string `json:"teacher_id"`
}

func GetTeaching(c *gin.Context) {
	var req GetTeachingReq
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

	var teaching model.ZhaoykTeaching
	if err := database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_id = ?", req.ID).Scan(&teaching).Error; err != nil || teaching.ID == 0 {
		apiException.AbortWithException(c, apiException.TeachingNotFoundError, err)
		return
	}

	resp := GetTeachingResp{
		ID:        teaching.ID,
		CourseID:  teaching.CourseID,
		TeacherID: teaching.TeacherID,
	}
	utils.JsonSuccessResponse(c, resp)
}

type GetTeacherListReq struct {
	AdminID int64 `form:"admin_id"`
}

type TeacherListItem struct {
	ID        int64  `json:"id"`
	TeacherID string `json:"teacher_id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	Title     string `json:"title"`
	Phone     string `json:"phone"`
	IsAdmin   bool   `json:"is_admin"`
}

type GetTeacherListResp struct {
	List []TeacherListItem `json:"list"`
}

func GetTeacherList(c *gin.Context) {
	var req GetTeacherListReq
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

	var teachers []model.ZhaoykTea
	if err := database.DB.Raw("SELECT * FROM zhaoyk_tea").Scan(&teachers).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]TeacherListItem, 0, len(teachers))
	for _, t := range teachers {
		list = append(list, TeacherListItem{
			ID:        t.ID,
			TeacherID: t.TeacherID,
			Name:      t.Name,
			Gender:    t.Gender,
			Age:       t.Age,
			Title:     t.Title,
			Phone:     t.Phone,
			IsAdmin:   t.IsAdmin,
		})
	}

	utils.JsonSuccessResponse(c, GetTeacherListResp{List: list})
}

type DeleteTeachingReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"id"`
}

func DeleteTeaching(c *gin.Context) {
	var req DeleteTeachingReq
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

	if err := database.DB.Exec("DELETE FROM zhaoyk_teaching WHERE zyk_id = ?", req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type CreateTeachingReq struct {
	AdminID   int64  `json:"admin_id"`
	TeacherID string `json:"teacher_id"`
	CourseID  int64  `json:"course_id"`
}

func CreateTeaching(c *gin.Context) {
	var req CreateTeachingReq
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

	if err := database.DB.Exec("INSERT INTO zhaoyk_teaching (zyk_teacher_id, zyk_course_id) VALUES (?, ?)", req.TeacherID, req.CourseID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type UpdateTeachingReq struct {
	AdminID   int64  `json:"admin_id"`
	ID        int64  `json:"id"`
	TeacherID string `json:"teacher_id"`
	CourseID  int64  `json:"course_id"`
}

func UpdateTeaching(c *gin.Context) {
	var req UpdateTeachingReq
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

	if err := database.DB.Exec("UPDATE zhaoyk_teaching SET zyk_course_id = ?, zyk_teacher_id = ? WHERE zyk_id = ?", req.CourseID, req.TeacherID, req.ID).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetTeachingListReq struct {
	AdminID int64 `form:"admin_id" binding:"required"` // 管理员ID
}

type TeachingListItem struct {
	ID        int64  `json:"id"`
	CourseID  int64  `json:"course_id"`
	TeacherID string `json:"teacher_id"`
}

func GetTeachingList(c *gin.Context) {
	var req GetTeachingListReq
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

	var teachings []model.ZhaoykTeaching
	if err := database.DB.Raw("SELECT * FROM zhaoyk_teaching").Scan(&teachings).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]TeachingListItem, len(teachings))
	for i, t := range teachings {
		list[i] = TeachingListItem{
			ID:        t.ID,
			CourseID:  t.CourseID,
			TeacherID: t.TeacherID,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"list": list,
	})
}
