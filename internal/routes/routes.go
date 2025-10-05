package routes

import (
	"testcase/config"
	"testcase/internal/infrastructures/database"
	"testcase/internal/middlewares"
	"testcase/internal/modules/document"
	documentHandler "testcase/internal/modules/document/handlers"
	documentRepository "testcase/internal/modules/document/repositories"
	documentService "testcase/internal/modules/document/services"
	"testcase/internal/modules/user"
	userHandler "testcase/internal/modules/user/handlers"
	userRepository "testcase/internal/modules/user/repositories"
	userService "testcase/internal/modules/user/services"
	"testcase/internal/utils"
	"testcase/package/securities"

	"github.com/gin-gonic/gin"
)

func InitHttpRoutes(r *gin.Engine, db *database.Database) {
	config := config.LoadConfig()
	jwtManager := securities.NewJWTManager(
		config.AccessTokenSecret,
		config.RefreshTokenSecret,
		config.TokenExpiry,
		config.RefreshExpiry,
		"",
	)

	authMware := middlewares.NewAuthMiddleware(jwtManager)

	userRepo := userRepository.NewUserRepository(db)
	documentRepo := documentRepository.NewDocumentRepository(db)

	userService := userService.NewUserService(userRepo, jwtManager)
	documentService := documentService.NewDocumentService(documentRepo)

	documentHandler := documentHandler.NewDocumentHandler(documentService)
	userHandler := userHandler.NewUserHandler(userService)

	v1 := r.Group("api/v1")
	{
		user.RegisterUserRoutes(v1, userHandler, authMware)
		document.RegisterDocumentRoutes(v1, documentHandler, authMware)
	}

	r.NoRoute(func(c *gin.Context) { utils.HandleRouteNotFound(c) })
}
