package adminController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetAvgScoreReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type courseAvgScoreView struct {
	CourseID   int64   `json:"course_id" gorm:"column:zyk_course_id"`
	CourseName string  `json:"course_name" gorm:"column:zyk_course_name"`
	SchoolYear int64   `json:"school_year" gorm:"column:zyk_school_year"`
	AvgScore   float64 `json:"avg_score" gorm:"column:zyk_avg_score"`
}

func GetAvgScore(c *gin.Context) {
	var req GetAvgScoreReq
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

	var avgScores []courseAvgScoreView
	if err := database.DB.Raw("SELECT * FROM zhaoyk_course_avg_score_view").Scan(&avgScores).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	result := make([]gin.H, len(avgScores))
	for i, row := range avgScores {
		result[i] = gin.H{
			"course_id":   row.CourseID,
			"course_name": row.CourseName,
			"school_year": row.SchoolYear,
			"avg_score":   row.AvgScore,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"avg_scores": result,
	})
}

type GetCityGpaTopTenCountReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type Top10GpaItem struct {
	ProvinceName   string `json:"province_name"`
	DepartmentName string `json:"department_name"`
	TopTenCount    int64  `json:"top10_count"`
}

type Top10GpaByProvinceDepartmentView struct {
	ProvinceName   string `gorm:"column:zyk_province_name"`
	DepartmentName string `gorm:"column:zyk_department_name"`
	Top10Count     int64  `gorm:"column:zyk_top10_count"`
}

func GetCityGpaTopTenCount(c *gin.Context) {
	var req GetCityGpaTopTenCountReq
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

	var top10Counts []Top10GpaByProvinceDepartmentView
	if err := database.DB.Raw("SELECT * FROM zhaoyk_top10_gpa_by_province_department_view").Scan(&top10Counts).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]gin.H, 0, len(top10Counts))
	for _, item := range top10Counts {
		list = append(list, gin.H{
			"province_name":   item.ProvinceName,
			"department_name": item.DepartmentName,
			"top10_count":     item.Top10Count,
		})
	}

	utils.JsonSuccessResponse(c, gin.H{
		"list": list,
	})
}

type GetCourseAllStudentScoreReq struct {
	AdminID    int64  `form:"admin_id"`
	CourseName string `form:"course_name"`
	SchoolYear int64  `form:"school_year"`
}

type StudentScoreItem struct {
	StudentID   int64    `json:"student_id"`
	StudentName string   `json:"student_name"`
	CourseName  string   `json:"course_name"`
	ClassName   string   `json:"class_name"`
	Semester    string   `json:"semester"`
	TeacherName []string `json:"teacher_name"`
	Score       float64  `json:"score"`
}

type GetCourseAllStudentScoreResp struct {
	Scores []StudentScoreItem `json:"scores"`
}

type CourseScore struct {
	StudentID   int64   `gorm:"column:zyk_student_id"`
	StudentName string  `gorm:"column:zyk_student_name"`
	CourseName  string  `gorm:"column:zyk_course_name"`
	ClassName   string  `gorm:"column:zyk_class_name"`
	Semester    string  `gorm:"column:zyk_semester"`
	TeacherName string  `gorm:"column:zyk_teacher_name"`
	Score       float64 `gorm:"column:zyk_score"`
}

func GetCourseAllStudentScore(c *gin.Context) {
	var req GetCourseAllStudentScoreReq
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

	var courseScores []CourseScore
	if err := database.DB.Raw(
		"SELECT * FROM zhaoyk_select_course_scores_by_name_and_year(?, ?)",
		req.CourseName, req.SchoolYear,
	).Scan(&courseScores).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	type groupKey struct {
		StudentID  int64
		CourseName string
		ClassName  string
		Semester   string
	}

	grouped := make(map[groupKey]*StudentScoreItem)
	for _, score := range courseScores {
		key := groupKey{
			StudentID:  score.StudentID,
			CourseName: score.CourseName,
			ClassName:  score.ClassName,
			Semester:   score.Semester,
		}

		if _, exists := grouped[key]; !exists {
			grouped[key] = &StudentScoreItem{
				StudentID:   score.StudentID,
				StudentName: score.StudentName,
				CourseName:  score.CourseName,
				ClassName:   score.ClassName,
				Semester:    score.Semester,
				Score:       score.Score,
				TeacherName: []string{},
			}
		}
		found := false
		for _, name := range grouped[key].TeacherName {
			if name == score.TeacherName {
				found = true
				break
			}
		}
		if !found {
			grouped[key].TeacherName = append(grouped[key].TeacherName, score.TeacherName)
		}
	}

	scoreList := make([]StudentScoreItem, 0, len(grouped))
	for _, item := range grouped {
		scoreList = append(scoreList, *item)
	}

	utils.JsonSuccessResponse(c, GetCourseAllStudentScoreResp{Scores: scoreList})
}

type GetGpaRankReq struct {
	AdminID int64 `form:"admin_id"`
	ID      int64 `form:"department_id"`
	Type    int64 `form:"type"`
}

type GPARankItem struct {
	Rank        int64   `json:"rank"`
	StudentID   string  `json:"student_id"`
	StudentName string  `json:"student_name"`
	ClassID     int64   `json:"class_id"`
	ClassName   string  `json:"class_name"`
	GPA         float64 `json:"gpa"`
}

type GetGPARankResp struct {
	Rank []GPARankItem `json:"rank"`
}

type DepartmentGpaRanking struct {
	Rank        int64   `gorm:"column:zyk_rank"`
	StudentID   string  `gorm:"column:zyk_student_id"`
	StudentName string  `gorm:"column:zyk_student_name"`
	ClassID     int64   `gorm:"column:zyk_class_id"`
	ClassName   string  `gorm:"column:zyk_class_name"`
	GPA         float64 `gorm:"column:zyk_gpa"`
}

func GetGpaRank(c *gin.Context) {
	var req GetGpaRankReq
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

	var rankings []DepartmentGpaRanking
	var err error
	if req.Type == 1 {
		err = database.DB.Raw("SELECT * FROM zhaoyk_select_department_gpa_ranking(?)", req.ID).Scan(&rankings).Error
	} else {
		err = database.DB.Raw("SELECT * FROM zhaoyk_select_class_gpa_ranking(?)", req.ID).Scan(&rankings).Error
	}
	if err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	rankList := make([]GPARankItem, len(rankings))
	for i, r := range rankings {
		rankList[i] = GPARankItem{
			Rank:        r.Rank,
			StudentID:   r.StudentID,
			StudentName: r.StudentName,
			ClassID:     r.ClassID,
			ClassName:   r.ClassName,
			GPA:         r.GPA,
		}
	}

	utils.JsonSuccessResponse(c, GetGPARankResp{
		Rank: rankList,
	})
}

type GetScoreListReq struct {
	AdminID int64 `form:"admin_id" binding:"required"`
}

type ScoreListItem struct {
	ID        int64   `json:"id"`
	StudentID string  `json:"student_id"`
	CourseID  int64   `json:"course_id"`
	Score     float64 `json:"score"`
}

type GetScoreListResp struct {
	List []ScoreListItem `json:"list"`
}

func GetScoreList(c *gin.Context) {
	var req GetScoreListReq
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

	var scores []model.ZhaoykScore
	if err := database.DB.Raw("SELECT * FROM zhaoyk_score").Scan(&scores).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	list := make([]ScoreListItem, 0, len(scores))
	for _, s := range scores {
		list = append(list, ScoreListItem{
			ID:        s.ID,
			StudentID: s.StudentID,
			CourseID:  s.CourseID,
			Score:     s.Score,
		})
	}

	utils.JsonSuccessResponse(c, GetScoreListResp{
		List: list,
	})
}
