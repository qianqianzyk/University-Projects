package router

import (
	"SGMS/app/controllers/adminController"
	"SGMS/app/controllers/studentController"
	"SGMS/app/controllers/teacherController"
	"SGMS/app/controllers/userController"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)

		student := api.Group("/student")
		{
			student.GET("/info", studentController.GetStudentInfo)
			student.PUT("/info", studentController.UpdateStudentInfo)
			student.GET("/class/course", studentController.GetClassCourse)
			student.POST("/course/select", studentController.SelectCourse)
			student.GET("/score", studentController.GetStudentScore)
			student.GET("/teacher/course", studentController.GetTeacherCourse)
		}

		teacher := api.Group("/teacher")
		{
			teacher.GET("/", teacherController.GetTeacher)
			teacher.GET("/course", teacherController.GetTeacherCourseByTeacher)
			teacher.GET("/course/avgscore", teacherController.GetCourseAvgScore)
			teacher.GET("/student", teacherController.GetStudentByTeacher)
			teacher.GET("/course/student", teacherController.GetCourseStudentList)
			teacher.POST("/course/student/score", teacherController.SetStudentScore)
		}

		admin := api.Group("/admin")
		{
			admin.GET("/avg/score", adminController.GetAvgScore)
			admin.GET("/cities", adminController.GetCityList)
			admin.GET("/city/gpa/top/ten/count", adminController.GetCityGpaTopTenCount)
			admin.GET("/city/student/count", adminController.GetCityStudentCount)
			admin.GET("/class", adminController.GetClass)
			admin.DELETE("/class", adminController.DeleteClass)
			admin.POST("/class", adminController.CreateClass)
			admin.PUT("/class", adminController.UpdateClass)
			admin.GET("/class/course", adminController.GetClassCourseList)
			admin.GET("/classes", adminController.GetClassList)
			admin.GET("/course", adminController.GetCourse)
			admin.DELETE("/course", adminController.DeleteCourse)
			admin.POST("/course", adminController.CreateCourse)
			admin.PUT("/course", adminController.UpdateCourse)
			admin.GET("/course/score/distribution", adminController.GetCourseScoreDistribution)
			admin.GET("/course/student/score", adminController.GetCourseAllStudentScore)
			admin.GET("/courses", adminController.GetCourseList)
			admin.GET("/department", adminController.GetDepartment)
			admin.DELETE("/department", adminController.DeleteDepartment)
			admin.POST("/department", adminController.CreateDepartment)
			admin.PUT("/department", adminController.UpdateDepartment)
			admin.GET("/department/gpa/rank", adminController.GetGpaRank)
			admin.GET("/departments", adminController.GetDepartmentList)
			admin.GET("/provinces", adminController.GetProvinceList)
			admin.GET("/scores", adminController.GetScoreList)
			admin.GET("/student", adminController.GetStudent)
			admin.DELETE("/student", adminController.DeleteStudent)
			admin.POST("/student", adminController.CreateStudent)
			admin.PUT("/student", adminController.UpdateStudent)
			admin.GET("/students", adminController.GetStudentList)
			admin.GET("/teacher", adminController.GetTeacher)
			admin.DELETE("/teacher", adminController.DeleteTeacher)
			admin.POST("/teacher", adminController.CreateTeacher)
			admin.PUT("/teacher", adminController.UpdateTeacher)
			admin.GET("/teacher/course", adminController.GetTeacherCourseList)
			admin.GET("/teacher/statistics", adminController.GetTeacherStatistics)
			admin.GET("/teaching", adminController.GetTeaching)
			admin.GET("/teachers", adminController.GetTeacherList)
			admin.DELETE("/teaching", adminController.DeleteTeaching)
			admin.POST("/teaching", adminController.CreateTeaching)
			admin.PUT("/teaching", adminController.UpdateTeaching)
			admin.GET("/teachings", adminController.GetTeachingList)
		}
	}
}
