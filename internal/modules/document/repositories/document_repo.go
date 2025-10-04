package repositories

import (
	"context"
	"testcase/internal/modules/document/entities"
)

type DocumentRepo interface {
	FindById(id string) (*entities.Document, error)
	CreateDocument(ctx context.Context, doc *entities.Document) error
	UpdateDocument(ctx context.Context, doc *entities.Document) error
}
