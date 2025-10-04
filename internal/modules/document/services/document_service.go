package services

import (
	"context"
	"testcase/internal/modules/document/dto"
	"testcase/internal/modules/document/entities"
)

type DocumentService interface {
	FindById(ctx context.Context, id string) (*entities.Document, error)
	CreateDocument(ctx context.Context, action *dto.CreateDocumentDTO) (*entities.Document, error)
	SubmitAction(ctx context.Context, id string, action *dto.UpdateDocumentDTO) (*entities.Document, error)
	ResubmitAction(ctx context.Context, id string) (*entities.Document, error)
}
