package repositories

import (
	"context"
	"testcase/internal/helpers"
	"testcase/internal/modules/user/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, user *entities.User) error
	ListUsers(ctx context.Context, params *helpers.PaginationParams) ([]entities.User, int64, error)
	DeleteUser(ctx context.Context, user *entities.User) error
}
