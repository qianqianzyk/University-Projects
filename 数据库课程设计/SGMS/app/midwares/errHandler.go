package midwares

import (
	"SGMS/app/apiException"
	"SGMS/app/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrHandler 中间件用于处理请求错误
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *apiException.Error

				ok := errors.As(err, &apiErr)

				if !ok {
					apiErr = apiException.ServerError
				}

				utils.JsonErrorResponse(c, apiErr.Code, apiErr.Msg)
				return
			}
		}
	}
}

// HandleNotFound 处理 404 错误
func HandleNotFound(c *gin.Context) {
	err := apiException.NotFound
	utils.JsonResponse(c, http.StatusNotFound, err.Code, err.Msg, nil)
}
