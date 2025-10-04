package repositories

import (
	"context"
	"errors"
	"fmt"
	"testcase/internal/infrastructures/database"
	"testcase/internal/modules/document/entities"

	"gorm.io/gorm"
)

type documentRepositoryImpl struct {
	db *database.Database
}

func NewDocumentRepository(db *database.Database) DocumentRepo {
	return &documentRepositoryImpl{
		db: db,
	}
}

func (r *documentRepositoryImpl) FindById(id string) (*entities.Document, error) {
	var doc entities.Document

	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("document with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find document by ID: %w", err)
	}

	return &doc, nil
}

func (r *documentRepositoryImpl) CreateDocument(ctx context.Context, doc *entities.Document) error {
	err := r.db.WithContext(ctx).Create(doc).Error
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

func (r *documentRepositoryImpl) UpdateDocument(ctx context.Context, doc *entities.Document) error {
	err := r.db.WithContext(ctx).Save(doc).Error
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}
