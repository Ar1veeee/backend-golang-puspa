package service

import (
	"backend-golang/internal/auth/dto"
	authErrors "backend-golang/internal/auth/errors"
	"backend-golang/internal/auth/repository"
	userErrors "backend-golang/internal/user/errors"
	userService "backend-golang/internal/user/service"
	"backend-golang/shared/constants"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req *dto.RegisterRequest) error
	LoginUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshTokenUser(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	LogoutUser(ctx context.Context, refreshToken string) error
}

type authService struct {
	authRepo    repository.AuthRepository
	userService userService.UserService
	validator   *authValidator
	mapper      *authMapper
}

func NewAuthService(authRepo repository.AuthRepository, userService userService.UserService) AuthService {
	return &authService{
		authRepo:    authRepo,
		userService: userService,
		validator:   newAuthValidator(),
		mapper:      newAuthMapper(),
	}
}

func (a *authService) RegisterUser(ctx context.Context, req *dto.RegisterRequest) error {
	if err := a.validator.validateRegisterRequest(req); err != nil {
		return err
	}

	exists, err := a.authRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return userErrors.ErrEmailExists
	}

	exists, err = a.authRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return userErrors.ErrUsernameExists
	}

	user, err := a.mapper.createRequestToRegister(req)
	if err != nil {
		return err
	}

	if err := a.authRepo.Register(ctx, user); err != nil {
		return fmt.Errorf("%w: %v", userErrors.ErrUserCreationFailed, err)
	}

	return nil
}

func (a *authService) LoginUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := a.validator.validateLoginRequest(req); err != nil {
		return nil, err
	}

	var user models.User

	if err := a.authRepo.Login(ctx, req.Username, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("username", req.Username).
				Msg("Login attempt failed: user not found")
			return nil, userErrors.ErrInvalidCredentials
		}
		log.Error().
			Str("username", req.Username).
			Msg("Failed to retrieve user from repository")
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserRetrievalFailed, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Warn().Str("username", req.Username).Msg("Login attempt failed: invalid password")
		return nil, userErrors.ErrInvalidCredentials
	}

	accessToken, err := helpers.GenerateToken(user.Id, constants.Role(user.Role))
	if err != nil {
		return nil, authErrors.ErrGenerateToken
	}

	refreshToken, err := a.createAndSaveRefreshToken(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("userId", user.Id).
		Str("username", user.Username).
		Msg("Successfully logged in")

	response := a.mapper.loginToResponse(&user)
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	return response, nil
}

func (a *authService) createAndSaveRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshTokenString, expiresAt, err := helpers.GenerateRefreshToken()
	if err != nil {
		return "", authErrors.ErrGenerateRefreshToken
	}

	refreshTokenModel := &models.RefreshToken{
		Id:        uuid.New().String(),
		UserId:    userID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := a.authRepo.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to save refresh token to database")
		return "", authErrors.ErrSaveRefreshToken
	}

	return refreshTokenString, nil
}

func (a *authService) RefreshTokenUser(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	var token models.RefreshToken

	if err := a.authRepo.RefreshToken(ctx, req.RefreshToken, &token); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, authErrors.ErrInvalidRefreshToken
		}
		return nil, err
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, authErrors.ErrTokenExpired
	}

	userResponse, err := a.userService.GetUserByID(ctx, token.UserId)
	if err != nil {
		log.Error().Err(err).Str("user_id", token.UserId).Msg("User associated with refresh token not found")
		return nil, userErrors.ErrUserNotFound
	}

	newAccessToken, err := helpers.GenerateToken(token.UserId, constants.Role(userResponse.Role))
	if err != nil {
		log.Error().Err(err).Str("user_id", token.UserId).Msg("Failed to generate new access token during refresh")
		return nil, authErrors.ErrGenerateToken
	}

	response := a.mapper.refreshTokeToResponse(&token)
	response.AccessToken = newAccessToken
	response.ExpiresAt = token.ExpiresAt.Format("2006-01-02 15:04:05")

	log.Info().Str("user_id", token.UserId).Msg("Access token refreshed successfully")
	return response, nil
}

func (a *authService) LogoutUser(ctx context.Context, refreshToken string) error {
	if err := a.authRepo.Logout(ctx, refreshToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msg("Logout attempt for a non-existent or already revoked token")
			return authErrors.ErrInvalidRefreshToken
		}
		log.Error().Err(err).Msg("Error during logout process")
		return err
	}
	log.Info().Msg("User logged out successfully")
	return nil
}

type authValidator struct{}

func newAuthValidator() *authValidator {
	return &authValidator{}
}

func (v *authValidator) validateRegisterRequest(req *dto.RegisterRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return helpers.IsValidPassword(req.Password)
}

func (v *authValidator) validateLoginRequest(req *dto.LoginRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}
	return nil
}

type authMapper struct{}

func newAuthMapper() *authMapper {
	return &authMapper{}
}

func (m *authMapper) createRequestToRegister(req *dto.RegisterRequest) (*models.User, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      string(constants.RoleUser),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *authMapper) loginToResponse(user *models.User) *dto.LoginResponse {
	return &dto.LoginResponse{
		Id:        user.Id,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		TokenType: "Bearer",
	}
}

func (m *authMapper) refreshTokeToResponse(token *models.RefreshToken) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		RefreshToken: token.Token,
		TokenType:    "Bearer",
	}
}
