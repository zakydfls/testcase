package handlers

import (
	"context"
	"net/http"
	"testcase/internal/middlewares"
	"testcase/internal/modules/document/dto"
	"testcase/internal/modules/document/services"
	"testcase/internal/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService services.DocumentService
}

func NewDocumentHandler(documentService services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var input dto.CreateDocumentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		middlewares.ValidationErrorResponse(c, err)
		return
	}

	role, _ := c.Get(utils.RoleContextKey)
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, utils.RoleContextKey, role)

	document, err := h.documentService.CreateDocument(ctx, &input)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, document, "Document created successfully", http.StatusCreated)
}

func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")

	document, err := h.documentService.FindById(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, document, "Document retrieved successfully", http.StatusOK)
}

func (h *DocumentHandler) SubmitAction(c *gin.Context) {
	var input dto.UpdateDocumentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		middlewares.ValidationErrorResponse(c, err)
		return
	}

	id := c.Param("id")

	ctx := c.Request.Context()

	document, err := h.documentService.SubmitAction(ctx, id, &input)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, document, "Action submitted successfully", http.StatusOK)
}

func (h *DocumentHandler) ResubmitAction(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	document, err := h.documentService.ResubmitAction(ctx, id)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, document, "Document resubmitted successfully", http.StatusOK)
}
