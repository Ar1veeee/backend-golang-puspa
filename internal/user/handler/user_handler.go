package handler

import (
	"backend-golang/internal/user/dto"
	userErrors "backend-golang/internal/user/errors"
	"backend-golang/internal/user/service"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) handleErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.JSON(http.StatusUnprocessableEntity, types.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	statusCode := http.StatusInternalServerError
	message := "An internal server errors occurred."

	switch {
	case errors.Is(err, userErrors.ErrUserNotFound):
		statusCode = http.StatusNotFound
		message = "User not found."
	case errors.Is(err, userErrors.ErrEmailExists):
		statusCode = http.StatusConflict
		message = "A user with that email already exists."
	case errors.Is(err, userErrors.ErrUsernameExists):
		statusCode = http.StatusConflict
		message = "A user with that username already exists."
	case errors.Is(err, userErrors.ErrInvalidUserID), errors.Is(err, userErrors.ErrUserIDRequired):
		statusCode = http.StatusBadRequest
		message = "Invalid or missing user ID."
	case errors.Is(err, globalErrors.ErrBadRequest), errors.Is(err, globalErrors.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = "Bad request."
	}

	c.JSON(statusCode, types.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  helpers.TranslateErrorMessage(err),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "User Created Successfully",
		Data:    user,
	})
}

func (h *UserHandler) FindUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List Data Users",
		Data:    users,
	})
}

func (h *UserHandler) FindUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "User Updated Successfully",
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "User Deleted Successfully",
	})
}
