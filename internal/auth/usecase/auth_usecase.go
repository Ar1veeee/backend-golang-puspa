package usecase

import (
	"backend-golang/internal/auth/delivery/http/dto"
	authErrors "backend-golang/internal/auth/errors"
	"backend-golang/internal/auth/repository"
	"backend-golang/shared/constants"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 15 * time.Minute
)

type AuthUseCase interface {
	RegisterService(ctx context.Context, req *dto.RegisterRequest) error
	VerifyEmailService(ctx context.Context, req *dto.VerifyCodeRequest) error
	ForgetPasswordService(ctx context.Context, req *dto.ForgetPasswordRequest) error
	LoginService(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshTokenService(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	LogoutService(ctx context.Context, refreshToken string) error
}

type authUseCase struct {
	authRepo    repository.AuthRepository
	validator   AuthValidator
	mapper      AuthMapper
	redisClient *redis.Client
}

func NewAuthUseCase(
	authRepo repository.AuthRepository,
	redisClient *redis.Client,
) AuthUseCase {
	return &authUseCase{
		authRepo:    authRepo,
		validator:   NewAuthValidator(),
		mapper:      NewAuthMapper(),
		redisClient: redisClient,
	}
}

func (uc *authUseCase) RegisterService(ctx context.Context, req *dto.RegisterRequest) error {
	if err := uc.validator.ValidateRegisterRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	exists, err := uc.authRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to check email existence")
		return globalErrors.ErrInternalServer
	}
	if exists {
		log.Warn().Str("email", req.Email).Msg("Registration failed: email already exists")
		return globalErrors.ErrEmailExists
	}

	exists, err = uc.authRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("Failed to check username existence")
		return globalErrors.ErrInternalServer
	}
	if exists {
		log.Warn().Str("username", req.Username).Msg("Registration failed: username already exists")
		return globalErrors.ErrUsernameExists
	}

	parent, err := uc.authRepo.GetParentByTempEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get parent of user")
		return globalErrors.ErrInternalServer
	}
	if parent == nil {
		log.Warn().Str("email", req.Email).Msg("No parent found with temp_email")
		return fmt.Errorf("no parent found with temp_email: %s", req.Email)
	}

	user, err := uc.mapper.RegisterRequestToUser(req)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user entity")
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.authRepo.CreateUser(ctx, user); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user in database")
		return globalErrors.ErrInternalServer
	}

	verificationCode, err := uc.mapper.CreateVerificationCodeWithEmail(user.Id, user.Email, user.Username)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create verification code")
	}

	if verificationCode != nil {
		if err := uc.authRepo.SaveVerificationCode(ctx, verificationCode); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
			return globalErrors.ErrInternalServer
		}
	}

	if err := uc.authRepo.UpdateParentUserId(ctx, parent.TempEmail, user.Id); err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to update parent of user")
		return globalErrors.ErrInternalServer
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("User registered and verification email sent")

	return nil
}

func (uc *authUseCase) VerifyEmailService(ctx context.Context, req *dto.VerifyCodeRequest) error {
	if err := uc.validator.ValidateVerifyCodeRequest(req); err != nil {
		log.Warn().Err(err).Str("code", req.Code).Msg("Verify code validation failed")
		return err
	}

	verificationCode, err := uc.authRepo.VerifyEmailByCode(ctx, req.Code)
	if err != nil {
		log.Warn().Err(err).Str("code", req.Code).Msg("Invalid verification code")
		return authErrors.ErrInvalidCode
	}

	if verificationCode.IsExpired() {
		log.Warn().Str("code", req.Code).Str("userId", verificationCode.UserId).Msg("Verification code expired")
	}

	if err := uc.authRepo.UpdateUserActiveStatus(ctx, verificationCode.UserId, true); err != nil {
		log.Error().Err(err).Str("userId", verificationCode.UserId).Msg("Failed to activate user")
		return globalErrors.ErrInternalServer
	}

	log.Info().
		Str("userId", verificationCode.UserId).
		Msg("Email verified successfully")

	return nil
}

func (uc *authUseCase) ForgetPasswordService(ctx context.Context, req *dto.ForgetPasswordRequest) error {
	if err := uc.validator.ValidateForgetPasswordRequest(req); err != nil {
		return err
	}

	user, err := uc.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Reset password attempt failed: email not found")
		return nil
	}

	verificationCode, err := uc.mapper.CreateForgerPasswordCode(user.Id, user.Email, user.Username)
	if err != nil {
		log.Warn().Err(err).Str("userId", user.Id).Msg("Failed to create forget password code")
		return globalErrors.ErrInternalServer
	}

	if err := uc.authRepo.SaveVerificationCode(ctx, verificationCode); err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
		return globalErrors.ErrInternalServer
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Forget password code sent")
	return nil
}

func (uc *authUseCase) LoginService(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := uc.validator.ValidateLoginRequest(req); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login validation failed")
		return nil, err
	}

	if err := uc.checkLoginRateLimit(ctx, req.Identifier); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login blocked due to rate limiting")
		return nil, err
	}

	user, err := uc.authRepo.FindUserByUsernameAndEmail(ctx, req.Identifier)
	if err != nil {
		uc.incrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("identifier", req.Identifier).Msg("Login attempt failed: user not found")
		return nil, globalErrors.ErrInvalidCredentials
	}

	if !user.IsActive {
		log.Warn().
			Str("identifier", req.Identifier).
			Msg("Login attempt failed: account is not active")
		return nil, globalErrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		uc.incrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("username", req.Identifier).Str("userId", user.Id).Msg("Login attempt failed: invalid password")
		return nil, globalErrors.ErrInvalidCredentials
	}

	uc.clearFailedAttempts(ctx, req.Identifier)

	accessToken, err := helpers.GenerateToken(user.Id, constants.Role(user.Role))
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to generate access token")
		return nil, authErrors.ErrGenerateToken
	}

	refreshToken, err := uc.createAndSaveRefreshToken(ctx, user.Id)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to generate refresh token")
		return nil, err
	}

	response := uc.mapper.UserToLoginResponse(user)
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	log.Info().
		Str("userId", user.Id).
		Str("username", user.Username).
		Msg("Successfully logged in")

	return response, nil
}

func (uc *authUseCase) RefreshTokenService(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		log.Warn().Msg("Refresh token cannot be nil")
		return nil, authErrors.ErrInvalidRefreshToken
	}

	refreshToken, err := uc.authRepo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		log.Warn().Str("token", req.RefreshToken[:10]+"...").Msg("Invalid refresh token provided")
		return nil, authErrors.ErrInvalidRefreshToken
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, authErrors.ErrInvalidRefreshToken
	}

	user, err := uc.authRepo.GetUserById(ctx, refreshToken.UserId)
	if err != nil {
		log.Error().Err(err).Str("user_id", refreshToken.UserId).Msg("User associated with refresh token not found")
		return nil, globalErrors.ErrUserNotFound
	}

	if !user.IsActive {
		log.Warn().Str("userId", user.Id).Msg("Refresh token request for inactive user")
		return nil, globalErrors.ErrUserInactive
	}

	newAccessToken, err := helpers.GenerateToken(refreshToken.UserId, constants.Role(user.Role))
	if err != nil {
		log.Error().Err(err).Str("user_id", refreshToken.UserId).Msg("Failed to generate new access token during refresh")
		return nil, globalErrors.ErrInternalServer
	}

	response := uc.mapper.RefreshTokenToResponse(refreshToken)
	response.AccessToken = newAccessToken
	response.ExpiresAt = refreshToken.ExpiresAt.Format("2006-01-02 15:04:05")

	log.Info().Str("user_id", refreshToken.UserId).Msg("Access token refreshed successfully")
	return response, nil
}

func (uc *authUseCase) LogoutService(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		log.Warn().Msg("Logout attempt with empty refresh token")
		return authErrors.ErrInvalidRefreshToken
	}

	if err := uc.authRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		log.Warn().Str("token", refreshToken[:10]+"...").Msg("Failed to revoke refresh token during logout")
		return authErrors.ErrInvalidRefreshToken
	}

	log.Info().Str("token", refreshToken[:10]+"...").Msg("User logged out successfully")
	return nil
}

func (uc *authUseCase) checkLoginRateLimit(ctx context.Context, identifier string) error {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)

	failedAttempts, err := uc.redisClient.Get(ctx, redisKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error().Err(err).Str("identifier", identifier).Msg("Failed to get login attempts from Redis")
		return globalErrors.ErrInternalServer
	}

	if failedAttempts >= maxFailedAttempts {
		log.Warn().Str("identifier", identifier).Int("attempts", failedAttempts).Msg("Login rate limit exceeded")
		return authErrors.ErrTooManyLoginAttempts
	}

	return nil
}

func (uc *authUseCase) incrementFailedAttempts(ctx context.Context, identifier string) {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)

	result := uc.redisClient.Incr(ctx, redisKey)
	if result.Err() != nil {
		log.Error().Err(result.Err()).Str("identifier", identifier).Msg("Failed to increment failed login attempts")
		return
	}

	attempts := result.Val()
	uc.redisClient.Expire(ctx, redisKey, lockoutDuration)

	log.Warn().Str("identifier", identifier).Int64("attempts", attempts).Msg("Failed login attempt recorded")
}

func (uc *authUseCase) clearFailedAttempts(ctx context.Context, identifier string) {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)

	if err := uc.redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Error().Err(err).Str("identifier", identifier).Msg("Failed to clear failed login attempts")
	}
}

func (uc *authUseCase) createAndSaveRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshToken := uc.mapper.CreateRefreshToken(userID)

	// Save using repository
	if err := uc.authRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		log.Error().Err(err).Str("userId", userID).Msg("Failed to save refresh token")
		return "", globalErrors.ErrInternalServer
	}

	return refreshToken.Token, nil
}
