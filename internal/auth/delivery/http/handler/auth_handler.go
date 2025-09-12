package handler

import (
	authDto "backend-golang/internal/auth/delivery/http/dto"
	"backend-golang/internal/auth/usecase"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService usecase.AuthUseCase
}

func NewAuthHandler(authService usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := authDto.RegisterRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.authService.RegisterUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "User created",
		Data:    nil,
	})
}

func (h *AuthHandler) ResendVerificationAccount(c *gin.Context) {
	email := c.Query("email")

	req := authDto.ResendTokenRequest{
		Email: email,
	}

	err := h.authService.ResendVerificationAccountUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "Email Verification Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) VerificationAccount(c *gin.Context) {
	token := c.Query("token")

	req := authDto.VerifyTokenRequest{
		Token: token,
	}

	err := h.authService.VerificationAccountUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Verification Account Successfully",
		Data:    nil,
	})
}

func (h *AuthHandler) ForgetPassword(c *gin.Context) {
	req := authDto.ForgetPasswordRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.authService.ForgetPasswordUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Email Successfully Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) ResendForgetPassword(c *gin.Context) {
	email := c.Query("email")

	req := authDto.ResendTokenRequest{
		Email: email,
	}

	err := h.authService.ResendForgetPasswordUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "Email Reset Password Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	token := c.Query("token")

	req := authDto.ResetPasswordRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	req.Token = token

	err := h.authService.ResetPasswordUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Password Reset Successfully",
		Data:    nil,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	req := authDto.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	userLogin, err := h.authService.LoginUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
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
		middlewares.AbortWithError(c, err)
		return
	}

	var refreshToken, err = h.authService.RefreshTokenUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
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
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.authService.LogoutUseCase(c.Request.Context(), refreshToken); err != nil {
		middlewares.AbortWithError(c, err)
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
