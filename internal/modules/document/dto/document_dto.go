package dto

import "testcase/internal/modules/document/entities"

type CreateDocumentDTO struct {
	Title string `json:"title" binding:"required"`
}

type UpdateDocumentDTO struct {
	Action  entities.DocumentAction `json:"action" binding:"required,oneof=approve reject"`
	Comment *string                 `json:"comment"`
}
