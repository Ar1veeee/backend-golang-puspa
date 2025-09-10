package service

import (
	"backend-golang/internal/auth/delivery/http/dto"
	authErrors "backend-golang/internal/auth/errors"
	"backend-golang/internal/auth/repository"
	"backend-golang/shared/constants"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 15 * time.Minute
)

type AuthService interface {
	RegisterService(ctx context.Context, req *dto.RegisterRequest) error
	VerifyEmailService(ctx context.Context, req *dto.VerifyCodeRequest) error
	ForgetPasswordService(ctx context.Context, req *dto.ForgetPasswordRequest) error
	LoginService(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshTokenService(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	LogoutService(ctx context.Context, refreshToken string) error
}

type authService struct {
	authRepo    repository.AuthRepository
	validator   *authValidator
	mapper      *authMapper
	redisClient *redis.Client
}

func NewAuthService(
	authRepo repository.AuthRepository,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		authRepo:    authRepo,
		validator:   newAuthValidator(),
		mapper:      newAuthMapper(),
		redisClient: redisClient,
	}
}

func (a *authService) RegisterService(ctx context.Context, req *dto.RegisterRequest) error {
	if err := a.validator.validateRegisterRequest(req); err != nil {
		return err
	}

	exists, err := a.authRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrEmailExists
	}

	exists, err = a.authRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrUsernameExists
	}

	parent, err := a.authRepo.GetParentByTempEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("failed to validate parent email: %w", err)
	}
	if parent == nil {
		return fmt.Errorf("no parent found with temp_email: %s", req.Email)
	}

	user, verificationCode, err := a.mapper.createRequestToRegister(req)
	if err != nil {
		return err
	}

	if err := a.authRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("%w: %v", authErrors.ErrUserCreationFailed, err)
	}

	if err := a.authRepo.SaveVerificationCode(ctx, verificationCode); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	parent.UserId = &user.Id
	if err := a.authRepo.InsertParentUserId(ctx, parent); err != nil {
		return fmt.Errorf("failed to update parent with user_id: %w", err)
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("User registered and verification email sent")

	return nil
}

func (a *authService) VerifyEmailService(ctx context.Context, req *dto.VerifyCodeRequest) error {
	if err := a.validator.validateVerifyEmailRequest(req); err != nil {
		return err
	}

	var verifyCode models.VerificationCode

	if err := a.authRepo.VerifyEmailByCode(ctx, req.Code, &verifyCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("identifier", req.Code).
				Msg("Verify attempt failed: code not found")
			return authErrors.ErrInvalidCode
		}
		if err.Error() == "verification code has expired" {
			log.Warn().
				Str("code", verifyCode.Code).
				Msg("Verify attempt failed: code has expired")
			return authErrors.ErrInvalidCode
		}
		log.Error().Err(err).Str("code", req.Code).Msg("Failed to verify email")
		return fmt.Errorf("verification failed: %w", err)
	}

	log.Info().
		Str("userId", verifyCode.UserId).
		Str("email", verifyCode.User.Email).
		Msg("Email verified successfully")

	return nil
}

func (a *authService) ForgetPasswordService(ctx context.Context, req *dto.ForgetPasswordRequest) error {
	if err := a.validator.validateForgetPasswordRequest(req); err != nil {
		return err
	}

	var user models.User
	if err := a.authRepo.FindUserByEmail(ctx, req.Email, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("email", req.Email).
				Msg("Reset password attempt failed: email not found")
			return globalErrors.ErrInvalidCredentials
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to reset password")
		return fmt.Errorf("reset password failed: %w", err)
	}

	verificationCode, err := a.mapper.createRequestToForgetPassword(req, a.authRepo)
	if err != nil {
		return err
	}

	if err := a.authRepo.SaveVerificationCode(ctx, verificationCode); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	return nil
}

func (a *authService) LoginService(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := a.validator.validateLoginRequest(req); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("login_attempts:%s", req.Identifier)

	failedAttempts, err := a.redisClient.Get(ctx, redisKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error().Err(err).Str("username", req.Identifier).Msg("Failed to get login attempts from Redis")
		return nil, globalErrors.ErrInternalServer
	}

	if failedAttempts >= maxFailedAttempts {
		log.Warn().Str("username", req.Identifier).Msg("Login attempt blocked due to too many failed attempts")
		return nil, authErrors.ErrTooManyLoginAttempts
	}

	var user models.User

	if err := a.authRepo.FindUserByUsernamAndEmail(ctx, req.Identifier, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.redisClient.Incr(ctx, redisKey)
			a.redisClient.Expire(ctx, redisKey, lockoutDuration)
			log.Warn().
				Str("identifier", req.Identifier).
				Msg("Login attempt failed: user not found")
			return nil, globalErrors.ErrInvalidCredentials
		}
		log.Error().
			Str("identifier", req.Identifier).
			Msg("Failed to retrieve user from repository")
		return nil, fmt.Errorf("%w: %v", authErrors.ErrUserRetrievalFailed, err)
	}

	if !user.IsActive {
		log.Warn().
			Str("identifier", req.Identifier).
			Msg("Login attempt failed: account is not active")
		return nil, globalErrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		a.redisClient.Incr(ctx, redisKey)
		a.redisClient.Expire(ctx, redisKey, lockoutDuration)
		log.Warn().Str("username", req.Identifier).Msg("Login attempt failed: invalid password")
		return nil, globalErrors.ErrInvalidCredentials
	}

	if err := a.redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Error().Err(err).Str("username", req.Identifier).Msg("Failed to delete user from Redis")
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

func (a *authService) RefreshTokenService(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
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

	userResponse, err := a.authRepo.GetUserById(ctx, token.UserId)
	if err != nil {
		log.Error().Err(err).Str("user_id", token.UserId).Msg("User associated with refresh token not found")
		return nil, globalErrors.ErrUserNotFound
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

func (a *authService) LogoutService(ctx context.Context, refreshToken string) error {
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

func (v *authValidator) validateVerifyEmailRequest(req *dto.VerifyCodeRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.Code == "" {
		return errors.New("verification code cannot be empty")
	}

	return nil
}

func (v *authValidator) validateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}
	if strings.Contains(req.Email, "@") {
		if !helpers.IsValidEmail(req.Email) {
			return fmt.Errorf("invalid email format")
		}
	}
	return nil
}

func (v *authValidator) validateLoginRequest(req *dto.LoginRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}
	if strings.Contains(req.Identifier, "@") {
		if !helpers.IsValidEmail(req.Identifier) {
			return fmt.Errorf("invalid email format")
		}
	}
	return nil
}

type authMapper struct{}

func newAuthMapper() *authMapper {
	return &authMapper{}
}

func (m *authMapper) createRequestToRegister(req *dto.RegisterRequest) (*models.User, *models.VerificationCode, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userID := helpers.GenerateULID()

	user := &models.User{
		Id:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      string(constants.RoleUser),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	code, err := helpers.GenerateVerificationCode()
	if err != nil {
		return nil, nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	verificationCode := &models.VerificationCode{
		UserId:    user.Id,
		Code:      code,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := helpers.SendEmail(req.Email, req.Username, code, "verifikasi_email_registrasi", "Verifikasi Email Anda"); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to send verification email")
		return user, verificationCode, err
	}

	return user, verificationCode, nil
}

func (m *authMapper) createRequestToForgetPassword(req *dto.ForgetPasswordRequest, authRepo repository.AuthRepository) (*models.VerificationCode, error) {

	var user models.User
	if err := authRepo.FindUserByUsernamAndEmail(context.Background(), req.Email, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, globalErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	code, err := helpers.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	verificationCode := &models.VerificationCode{
		UserId:    user.Id,
		Code:      code,
		Status:    string(constants.VerificationCodeStatusPending),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := helpers.SendEmail(user.Email, user.Username, code, "forget_password_email", "Reset Password Anda"); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to send email")
		return nil, err
	}

	return verificationCode, nil
}

func (m *authMapper) loginToResponse(user *models.User) *dto.LoginResponse {
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

func (m *authMapper) refreshTokeToResponse(token *models.RefreshToken) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		RefreshToken: token.Token,
		TokenType:    "Bearer",
	}
}
