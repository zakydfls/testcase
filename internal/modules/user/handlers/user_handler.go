package handlers

import (
	"net/http"
	"testcase/internal/middlewares"
	"testcase/internal/modules/user/dto"
	"testcase/internal/modules/user/services"
	"testcase/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middlewares.ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, user, "User created successfully", http.StatusCreated)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var input dto.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middlewares.ValidationErrorResponse(c, err)
		return
	}

	loginResponse, err := h.userService.LoginUser(c.Request.Context(), &input)
	if err != nil {
		panic(err)
	}

	utils.SuccessResponse(c, loginResponse, "Login successful", http.StatusOK)
}
