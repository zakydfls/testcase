package document

import (
	"testcase/internal/middlewares"
	"testcase/internal/modules/document/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDocumentRoutes(rg *gin.RouterGroup, h *handlers.DocumentHandler, authMware *middlewares.AuthMiddleware) {

	documentRoutes := rg.Group("/documents")
	documentRoutes.Use(authMware.Auth())
	{
		documentRoutes.POST("/", h.CreateDocument)
		documentRoutes.POST("/:id/action", h.SubmitAction)
		documentRoutes.GET("/:id", h.GetDocument)
		documentRoutes.PUT("/:id", h.ResubmitAction)
	}
}
