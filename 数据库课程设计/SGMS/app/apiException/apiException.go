package apiException

import (
	"SGMS/app/utils/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error 表示自定义错误，包括状态码、消息和日志级别
type Error struct {
	Code  int
	Msg   string
	Level log.Level
}

// Error 表示自定义的错误类型
var (
	ServerError             = NewError(200500, log.LevelError, "系统异常，请稍后重试!")
	ParamError              = NewError(200501, log.LevelInfo, "参数错误")
	UserNotFoundError       = NewError(200502, log.LevelInfo, "用户不存在")
	ClassNotFoundError      = NewError(200503, log.LevelInfo, "班级不存在")
	CityNotFoundError       = NewError(200504, log.LevelInfo, "城市不存在")
	UpdateError             = NewError(200505, log.LevelInfo, "更新失败")
	DatabaseError           = NewError(200507, log.LevelInfo, "数据库调用失败")
	StudentNotInClassError  = NewError(200508, log.LevelInfo, "学生不属于该班级")
	TeacherNotFoundError    = NewError(200509, log.LevelInfo, "教师不存在")
	CourseNotFoundError     = NewError(200510, log.LevelInfo, "课程不存在")
	UsernameOrPasswordError = NewError(200511, log.LevelInfo, "用户名或密码错误")
	UserPermissionError     = NewError(200512, log.LevelInfo, "用户权限不足")
	UserRoleError           = NewError(200513, log.LevelInfo, "未知用户权限")
	DepartmentNotFoundError = NewError(200514, log.LevelInfo, "部门不存在")
	TeachingNotFoundError   = NewError(200515, log.LevelInfo, "授课信息不存在")

	NotFound = NewError(200404, log.LevelWarn, http.StatusText(http.StatusNotFound))
)

// Error 方法实现了 error 接口，返回错误的消息内容
func (e *Error) Error() string {
	return e.Msg
}

// NewError 创建并返回一个新的自定义错误实例
func NewError(code int, level log.Level, msg string) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiError *Error, err error) {
	fmt.Println(err)
	_ = c.AbortWithError(200, apiError)
}
