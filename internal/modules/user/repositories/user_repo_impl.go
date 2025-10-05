package repositories

import (
	"context"
	"errors"
	"fmt"

	"testcase/internal/helpers"
	"testcase/internal/infrastructures/database"
	"testcase/internal/modules/user/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	var user entities.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) ListUsers(ctx context.Context, params *helpers.PaginationParams) ([]entities.User, int64, error) {
	var users []entities.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.User{})

	if params.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if params.Filter != "" {
		query = query.Where("status = ?", params.Filter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	err := query.
		Offset(params.GetOffset()).
		Limit(params.Limit).
		Order(params.GetOrderBy()).
		Find(&users).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users with pagination: %w", err)
	}

	return users, total, nil
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Delete(user).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (*entities.User, error) {
	var user entities.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &user, nil
}
