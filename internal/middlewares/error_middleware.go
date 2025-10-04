package middlewares

import (
	"fmt"
	"testcase/internal/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *utils.AppError:
					utils.ErrorResponse(c, e.ErrorCode, e.Err)
				case error:
					utils.ErrorResponse(c, utils.ErrFetchDataError, e.Error())
				default:
					utils.ErrorResponse(c, utils.ErrInternalServer, fmt.Sprint(e))
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
