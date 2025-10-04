package user

import (
	"testcase/internal/modules/user/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *handlers.UserHandler) {

	userRoutes := rg.Group("/users")
	{
		userRoutes.POST("/", h.CreateUser)
		userRoutes.POST("/login", h.LoginUser)
	}
}
