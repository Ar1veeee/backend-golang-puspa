package handlers

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	RegisterUC                  auth.RegisterUseCase
	LoginUC                     auth.LoginUseCase
	ResetPasswordUC             auth.ResetPasswordUseCase
	RefreshTokenUC              auth.RefreshTokenUseCase
	LogoutUC                    auth.LogoutUseCase
	ResendVerificationAccountUC auth.ResendVerificationAccountUseCase
	VerificationAccountUC       auth.VerificationAccountUseCase
	ForgetPasswordUC            auth.ForgetPasswordUseCase
	ResendForgetPasswordUC      auth.ResendForgetPasswordUseCase
}

func NewAuthHandler(
	registerUC auth.RegisterUseCase,
	loginUC auth.LoginUseCase,
	resetPasswordUC auth.ResetPasswordUseCase,
	refreshTokenUC auth.RefreshTokenUseCase,
	logoutUC auth.LogoutUseCase,
	resendVerificationAccountUC auth.ResendVerificationAccountUseCase,
	verificationAccountUC auth.VerificationAccountUseCase,
	forgetPasswordUC auth.ForgetPasswordUseCase,
	resendForgetPasswordUC auth.ResendForgetPasswordUseCase,
) *AuthHandler {
	return &AuthHandler{
		RegisterUC:                  registerUC,
		LoginUC:                     loginUC,
		ResetPasswordUC:             resetPasswordUC,
		RefreshTokenUC:              refreshTokenUC,
		LogoutUC:                    logoutUC,
		ResendVerificationAccountUC: resendVerificationAccountUC,
		VerificationAccountUC:       verificationAccountUC,
		ForgetPasswordUC:            forgetPasswordUC,
		ResendForgetPasswordUC:      resendForgetPasswordUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	req := dto.RegisterRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.RegisterUC.Execute(c.Request.Context(), &req)
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

	req := dto.ResendTokenRequest{
		Email: email,
	}

	err := h.ResendVerificationAccountUC.Execute(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Email Verification Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) VerificationAccount(c *gin.Context) {
	token := c.Query("token")

	req := dto.VerifyTokenRequest{
		Token: token,
	}

	err := h.VerificationAccountUC.Execute(c.Request.Context(), &req)
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
	req := dto.ForgetPasswordRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.ForgetPasswordUC.Execute(c.Request.Context(), &req)
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

	req := dto.ResendTokenRequest{
		Email: email,
	}

	err := h.ResendForgetPasswordUC.Execute(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Email Reset Password Sent",
		Data:    nil,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	token := c.Query("token")

	req := dto.ResetPasswordRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	req.Token = token

	err := h.ResetPasswordUC.Execute(c.Request.Context(), &req)
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
	req := dto.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	userLogin, err := h.LoginUC.Execute(c.Request.Context(), &req)
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
	req := dto.RefreshTokenRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	var refreshToken, err = h.RefreshTokenUC.Execute(c.Request.Context(), &req)
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

	if err := h.LogoutUC.Execute(c.Request.Context(), refreshToken); err != nil {
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
