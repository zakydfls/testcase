package services

import (
	"context"
	"fmt"
	"testcase/internal/modules/user/dto"
	"testcase/internal/modules/user/entities"
	"testcase/internal/modules/user/repositories"
	"testcase/internal/modules/user/types"
	"testcase/internal/utils"
	"testcase/package/securities"
	"time"
)

type userServiceImpl struct {
	userRepo   repositories.UserRepository
	jwtManager *securities.JWTManager
}

func (u *userServiceImpl) CreateUser(ctx context.Context, input *dto.CreateUserInput) (*entities.User, error) {
	username, _ := u.userRepo.FindByUsername(input.Username)
	if username != nil {
		return nil, utils.NewAppError(utils.ErrUsernameExists, fmt.Errorf("username already exists"))
	}
	email, _ := u.userRepo.FindByEmail(input.Email)
	if email != nil {
		return nil, utils.NewAppError(utils.ErrEmailExists, fmt.Errorf("email already exists"))
	}

	hashedPassword, hashErr := securities.HashPassword(input.Password)
	if hashErr != nil {
		return nil, hashErr
	}
	input.Password = hashedPassword

	user := &entities.User{
		Name:     input.Name,
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
		Phone:    input.Phone,
		Role:     input.Role,
		IsActive: true,
	}

	createUserErr := u.userRepo.CreateUser(ctx, user)
	if createUserErr != nil {
		return nil, createUserErr
	}
	return user, nil
}

func (u *userServiceImpl) LoginUser(ctx context.Context, input *dto.LoginUserInput) (*types.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, utils.NewAppError(utils.ErrNotFound, fmt.Errorf("user with email %s not found", input.Email))
	}
	if !user.IsActive {
		return nil, utils.NewAppError(utils.ErrInactiveUser, fmt.Errorf("user with email %s is inactive", input.Email))
	}

	errPass := securities.VerifyPassword(user.Password, input.Password)
	if errPass != nil {
		return nil, utils.NewAppError(utils.ErrInvalidCredentials, fmt.Errorf("invalid credentials"))
	}

	accessToken, refreshToken, err := u.jwtManager.GenerateTokenPair(&securities.JWTPayload{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLogin = &now
	updateErr := u.userRepo.UpdateUser(ctx, user)
	if updateErr != nil {
		return nil, updateErr
	}

	return &types.LoginResponse{
		Token: securities.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: *user,
	}, nil
}

func NewUserService(userRepo repositories.UserRepository, jwtManager *securities.JWTManager) UserService {
	return &userServiceImpl{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}
