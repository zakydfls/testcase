package repositories

import (
	"context"
	"errors"
	"fmt"
	"testcase/internal/helpers"
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

func (r *documentRepositoryImpl) ListDocuments(ctx context.Context, params *helpers.PaginationParams) ([]entities.Document, int64, error) {
	var docs []entities.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Document{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	if params != nil {
		limit := params.Limit
		offset := (params.Page - 1) * params.Limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&docs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}

	return docs, total, nil
}
