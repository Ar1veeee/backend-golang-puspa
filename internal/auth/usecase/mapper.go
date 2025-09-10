package usecase

import (
	"backend-golang/internal/auth/delivery/http/dto"
	"backend-golang/internal/auth/entity"
	"backend-golang/shared/constants"
	"backend-golang/shared/helpers"
	"time"

	"github.com/rs/zerolog/log"
)

type AuthMapper interface {
	RegisterRequestToUser(req *dto.RegisterRequest) (*entity.User, error)
	CreateVerificationCode(userId string) *entity.VerificationCode
	CreateRefreshToken(userId string) *entity.RefreshToken
	UserToLoginResponse(user *entity.User) *dto.LoginResponse
	RefreshTokenToResponse(token *entity.RefreshToken) *dto.RefreshTokenResponse
	CreateVerificationCodeWithEmail(userId, email, username string) (*entity.VerificationCode, error)
	CreateForgerPasswordCode(userId, email, username string) (*entity.VerificationCode, error)
}

type authMapper struct{}

func NewAuthMapper() AuthMapper {
	return &authMapper{}
}

func (m *authMapper) RegisterRequestToUser(req *dto.RegisterRequest) (*entity.User, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userId := helpers.GenerateULID()

	user := &entity.User{
		Id:        userId,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      string(constants.RoleUser),
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

func (m *authMapper) CreateVerificationCode(userId string) *entity.VerificationCode {
	code, _ := helpers.GenerateVerificationCode()

	return &entity.VerificationCode{
		UserId:    userId,
		Code:      code,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *authMapper) CreateRefreshToken(userId string) *entity.RefreshToken {
	tokenString, expiresAt, _ := helpers.GenerateRefreshToken()

	return &entity.RefreshToken{
		UserId:    userId,
		Token:     tokenString,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *authMapper) UserToLoginResponse(user *entity.User) *dto.LoginResponse {
	return &dto.LoginResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "Bearer",
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *authMapper) RefreshTokenToResponse(token *entity.RefreshToken) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		RefreshToken: token.Token,
		TokenType:    "Bearer",
		ExpiresAt:    token.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *authMapper) CreateVerificationCodeWithEmail(userId, email, username string) (*entity.VerificationCode, error) {
	code, err := helpers.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	verificationCode := &entity.VerificationCode{
		UserId:    userId,
		Code:      code,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := helpers.SendEmail(email, username, code, "verification_email", "Verifikasi Email Anda"); err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to send verification email")
		return verificationCode, err
	}

	return verificationCode, nil
}

func (m *authMapper) CreateForgerPasswordCode(userId, email, username string) (*entity.VerificationCode, error) {
	code, err := helpers.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	verificationCode := &entity.VerificationCode{
		UserId:    userId,
		Code:      code,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := helpers.SendEmail(email, username, code, "forget_password_email", "Reset Password Anda"); err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to send forget password email")
		return verificationCode, err
	}

	return verificationCode, nil
}
