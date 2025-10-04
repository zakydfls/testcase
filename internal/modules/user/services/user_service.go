package services

import (
	"context"
	"testcase/internal/modules/user/dto"
	"testcase/internal/modules/user/entities"
	"testcase/internal/modules/user/types"
)

type UserService interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) (*entities.User, error)
	LoginUser(ctx context.Context, input *dto.LoginUserInput) (*types.LoginResponse, error)
}
