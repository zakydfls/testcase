package user

import (
	"testcase/internal/middlewares"
	"testcase/internal/modules/user/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *handlers.UserHandler, authMware *middlewares.AuthMiddleware) {

	userRoutes := rg.Group("/users")
	{
		userRoutes.POST("/", h.CreateUser)
		userRoutes.POST("/login", h.LoginUser)
		userRoutes.POST("/refresh-token", authMware.AuthRefresh(), h.RefreshToken)
	}
}
