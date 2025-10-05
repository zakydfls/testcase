package repositories

import (
	"context"
	"testcase/internal/helpers"
	"testcase/internal/modules/document/entities"
)

type DocumentRepo interface {
	FindById(id string) (*entities.Document, error)
	CreateDocument(ctx context.Context, doc *entities.Document) error
	UpdateDocument(ctx context.Context, doc *entities.Document) error
	ListDocuments(ctx context.Context, params *helpers.PaginationParams) ([]entities.Document, int64, error)
}
