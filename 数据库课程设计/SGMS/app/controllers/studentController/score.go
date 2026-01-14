package studentController

import (
	"SGMS/app/apiException"
	"SGMS/app/model"
	"SGMS/app/utils"
	"SGMS/config/database"
	"github.com/gin-gonic/gin"
)

type GetStudentScoreReq struct {
	ID   int64 `form:"id" binding:"required"`
	Year int64 `form:"year"`
}

type ScoreList struct {
	CourseID       int64         `json:"course_id"`
	CourseName     string        `json:"course_name"`
	SchoolYear     int64         `json:"school_year"`
	Semester       string        `json:"semester"`
	Hours          int64         `json:"hours"`
	ExamType       string        `json:"exam_type"`
	TeacherList    []TeacherList `json:"teacher_list"`
	Credits        float64       `json:"credit"`
	Score          float64       `json:"score"`
	RetakeRequired bool          `json:"retake_required"`
}

func GetStudentScore(c *gin.Context) {
	var req GetStudentScoreReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var student model.ZhaoykStu
	if err := database.DB.Raw("SELECT * FROM zhaoyk_stu WHERE zyk_id = ?", req.ID).Scan(&student).Error; err != nil || student.ID == 0 {
		apiException.AbortWithException(c, apiException.UserNotFoundError, err)
		return
	}

	var rawScores []struct {
		CourseID       int64   `gorm:"column:zyk_course_id"`
		CourseName     string  `gorm:"column:zyk_course_name"`
		Credit         float64 `gorm:"column:zyk_credit"`
		Score          float64 `gorm:"column:zyk_score"`
		RetakeRequired string  `gorm:"column:zyk_retake_required"`
	}

	if err := database.DB.Raw("SELECT * FROM Zhaoyk_select_student_scores_and_retake_status(?, ?)",
		student.StudentID, req.Year).Scan(&rawScores).Error; err != nil {
		apiException.AbortWithException(c, apiException.DatabaseError, err)
		return
	}

	scoreList := make([]ScoreList, 0)
	for _, score := range rawScores {
		var course model.ZhaoykCou
		if err := database.DB.Raw("SELECT * FROM zhaoyk_cou WHERE zyk_id = ?", score.CourseID).Scan(&course).Error; err != nil {
			apiException.AbortWithException(c, apiException.CourseNotFoundError, err)
			return
		}

		var teachings []model.ZhaoykTeaching
		database.DB.Raw("SELECT * FROM zhaoyk_teaching WHERE zyk_course_id = ?", course.ID).Scan(&teachings)

		teachers := make([]TeacherList, len(teachings))
		for i, t := range teachings {
			var teacher model.ZhaoykTea
			if err := database.DB.Raw("SELECT * FROM zhaoyk_tea WHERE zyk_teacher_id = ?", t.TeacherID).Scan(&teacher).Error; err != nil || teacher.ID == 0 {
				apiException.AbortWithException(c, apiException.TeacherNotFoundError, err)
				return
			}
			teachers[i] = TeacherList{
				TeacherID:   teacher.ID,
				TeacherName: teacher.Name,
			}
		}

		retake := false
		if score.RetakeRequired == "是" {
			retake = true
		}

		scoreList = append(scoreList, ScoreList{
			CourseID:       score.CourseID,
			CourseName:     score.CourseName,
			SchoolYear:     course.SchoolYear,
			Semester:       course.Semester,
			Hours:          course.Hours,
			ExamType:       course.ExamType,
			TeacherList:    teachers,
			Credits:        score.Credit,
			Score:          score.Score,
			RetakeRequired: retake,
		})
	}

	utils.JsonSuccessResponse(c, gin.H{
		"scores": scoreList,
	})
}
