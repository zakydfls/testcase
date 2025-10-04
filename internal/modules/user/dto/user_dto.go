package dto

import "testcase/internal/modules/user/entities"

type CreateUserInput struct {
	Name     string            `json:"name" binding:"required"`
	Username string            `json:"username" binding:"required,alphanum,min=3,max=100"`
	Email    string            `json:"email" binding:"required,email"`
	Password string            `json:"password" binding:"required,min=6"`
	Phone    string            `json:"phone,omitempty"`
	Role     entities.RoleEnum `json:"role" binding:"required"`
}

type UpdateUserInput struct {
	Name     *string            `json:"name,omitempty"`
	Username *string            `json:"username,omitempty" binding:"omitempty,alphanum,min=3,max=100"`
	Email    *string            `json:"email,omitempty" binding:"omitempty,email"`
	Password *string            `json:"password,omitempty" binding:"omitempty,min=6"`
	Phone    *string            `json:"phone,omitempty"`
	Role     *entities.RoleEnum `json:"role,omitempty"`
	IsActive *bool              `json:"is_active,omitempty"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
