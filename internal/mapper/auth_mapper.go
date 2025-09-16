package mapper

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/helpers"
	"time"
)

type AuthMapper interface {
	RegisterRequestToUser(req *dto.RegisterRequest) (*entities.User, error)

	CreateVerificationToken(userId string) (*entities.VerificationToken, error)
	CreateRefreshToken(userId string) (*entities.RefreshToken, error)

	ResetPasswordRequestToUser(req *dto.ResetPasswordRequest) (*entities.User, error)

	LoginResponse(user *entities.User, refreshToken *entities.RefreshToken) *dto.LoginResponse
	RefreshTokenToResponse(token *entities.RefreshToken) *dto.RefreshTokenResponse
}

type authMapper struct{}

func NewAuthMapper() AuthMapper {
	return &authMapper{}
}

func (m *authMapper) RegisterRequestToUser(req *dto.RegisterRequest) (*entities.User, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userId := helpers.GenerateULID()

	user := &entities.User{
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

func (m *authMapper) CreateVerificationToken(userId string) (*entities.VerificationToken, error) {
	token, expiresAt, err := helpers.GenerateVerificationToken(userId)
	if err != nil {
		return nil, err
	}

	verificationCode := &entities.VerificationToken{
		UserId:    userId,
		Token:     token,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return verificationCode, nil
}

func (m *authMapper) CreateRefreshToken(userId string) (*entities.RefreshToken, error) {
	tokenString, expiresAt, _ := helpers.GenerateRefreshToken()

	refreshToken := &entities.RefreshToken{
		UserId:    userId,
		Token:     tokenString,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return refreshToken, nil
}

func (m *authMapper) ResetPasswordRequestToUser(req *dto.ResetPasswordRequest) (*entities.User, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Password:  hashedPassword,
		UpdatedAt: time.Now(),
	}

	return user, nil
}

func (m *authMapper) LoginResponse(user *entities.User, refreshToken *entities.RefreshToken) *dto.LoginResponse {
	return &dto.LoginResponse{
		Id:           user.Id,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		TokenType:    "Bearer",
		RefreshToken: refreshToken.Token,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *authMapper) RefreshTokenToResponse(token *entities.RefreshToken) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		RefreshToken: token.Token,
		TokenType:    "Bearer",
		ExpiresAt:    token.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
}
