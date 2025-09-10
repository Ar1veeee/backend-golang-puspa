package handler

import (
	authDto "backend-golang/internal/auth/delivery/http/dto"
	authErrors "backend-golang/internal/auth/errors"
	"backend-golang/internal/auth/usecase"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService usecase.AuthUseCase
}

func NewAuthHandler(authService usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) handleErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.JSON(http.StatusUnprocessableEntity, types.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var statusCode int
	var message string

	switch {
	case errors.Is(err, globalErrors.ErrUserNotFound):
		statusCode = http.StatusNotFound
		message = "User not found."
	case errors.Is(err, globalErrors.ErrEmailExists):
		statusCode = http.StatusConflict
		message = "A user with that email already exists."
	case errors.Is(err, globalErrors.ErrUsernameExists):
		statusCode = http.StatusConflict
		message = "A user with that username already exists."
	case errors.Is(err, globalErrors.ErrInvalidUserID), errors.Is(err, globalErrors.ErrUserIDRequired):
		statusCode = http.StatusBadRequest
		message = "Invalid or missing user ID."
	case errors.Is(err, globalErrors.ErrBadRequest), errors.Is(err, globalErrors.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = "Bad request."
	case errors.Is(err, globalErrors.ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
		message = "Invalid credentials."
	case errors.Is(err, globalErrors.ErrUserInactive):
		statusCode = http.StatusForbidden
		message = "Account is not active. Please verify your email."
	case errors.Is(err, globalErrors.ErrInternalServer):
		statusCode = http.StatusInternalServerError
		message = "An internal server error occurred."
	case errors.Is(err, authErrors.ErrInvalidCode):
		statusCode = http.StatusBadRequest
		message = "Invalid or expired verification code."
	case errors.Is(err, authErrors.ErrGenerateToken):
		statusCode = http.StatusInternalServerError
		message = "Failed to generate authentication token."
	case errors.Is(err, authErrors.ErrTooManyLoginAttempts):
		statusCode = http.StatusTooManyRequests
		message = "Too many login attempts. Please try again later."
	case errors.Is(err, authErrors.ErrInvalidRefreshToken):
		statusCode = http.StatusUnauthorized
		message = "Invalid or expired refresh token."
	default:
		statusCode = http.StatusInternalServerError
		message = "An unexpected error occurred."
	}

	c.JSON(statusCode, types.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  map[string]string{"error": message},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := authDto.RegisterRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	err := h.authService.RegisterService(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "User created",
		Data:    nil,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	req := authDto.VerifyCodeRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	err := h.authService.VerifyEmailService(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Verify Email Successfully",
		Data:    nil,
	})
}

func (h *AuthHandler) ForgetPassword(c *gin.Context) {
	req := authDto.ForgetPasswordRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	err := h.authService.ForgetPasswordService(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Email Successfully Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	req := authDto.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	userLogin, err := h.authService.LoginService(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.SetCookie(
		"refresh_token",
		userLogin.RefreshToken,
		3600*24*7,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data:    userLogin,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	req := authDto.RefreshTokenRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	var refreshToken, err = h.authService.RefreshTokenService(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "New Access Token Generated Successfully",
		Data:    refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.handleErrorResponse(c, globalErrors.ErrUnauthorized)
		return
	}

	if err := h.authService.LogoutService(c.Request.Context(), refreshToken); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Logout Success",
		Data:    nil,
	})
}
