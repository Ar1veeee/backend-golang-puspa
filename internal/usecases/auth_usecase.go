package usecases

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/domain/services"
	"backend-golang/internal/errors"
	"backend-golang/internal/helpers"
	"backend-golang/internal/mapper"
	"backend-golang/internal/validator"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterUseCase(ctx context.Context, req *dto.RegisterRequest) error
	LoginUseCase(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	ResetPasswordUseCase(ctx context.Context, req *dto.ResetPasswordRequest) error
	RefreshTokenUseCase(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	LogoutUseCase(ctx context.Context, refreshToken string) error

	ResendVerificationAccountUseCase(ctx context.Context, req *dto.ResendTokenRequest) error
	VerificationAccountUseCase(ctx context.Context, req *dto.VerifyTokenRequest) error
	ForgetPasswordUseCase(ctx context.Context, req *dto.ForgetPasswordRequest) error
	ResendForgetPasswordUseCase(ctx context.Context, req *dto.ResendTokenRequest) error
}

type authUseCase struct {
	txRepo           repositories.TransactionRepository
	userRepo         repositories.UserRepository
	parentRepo       repositories.ParentRepository
	verifyTokenRepo  repositories.VerificationTokenRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	validator        validator.AuthValidator
	mapper           mapper.AuthMapper
	emailService     services.EmailService
	rateLimiter      services.RateLimiterService
	tokenService     services.TokenService
}

func NewAuthUseCase(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	parentRepo repositories.ParentRepository,
	verifyTokenRepo repositories.VerificationTokenRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	emailService services.EmailService,
	rateLimiter services.RateLimiterService,
	tokenService services.TokenService,
) AuthUseCase {
	return &authUseCase{
		txRepo:           txRepo,
		userRepo:         userRepo,
		parentRepo:       parentRepo,
		verifyTokenRepo:  verifyTokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		validator:        validator.NewAuthValidator(),
		mapper:           mapper.NewAuthMapper(),
		emailService:     emailService,
		rateLimiter:      rateLimiter,
		tokenService:     tokenService,
	}
}

func (uc *authUseCase) RegisterUseCase(ctx context.Context, req *dto.RegisterRequest) error {
	if err := uc.validator.ValidateRegisterRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	emailExists, usernameExists, err := uc.userRepo.CheckExisting(ctx, req.Email, req.Username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if emailExists {
		tx.Rollback()
		return errors.ErrEmailExists
	}

	if usernameExists {
		tx.Rollback()
		return errors.ErrUsernameExists
	}

	parent, err := uc.parentRepo.GetByTempEmail(ctx, req.Email)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get parent of user")
		return errors.ErrInternalServer
	}
	if parent == nil {
		tx.Rollback()
		log.Warn().Str("email", req.Email).Msg("No parent found with temp_email")
		return errors.ErrEmailNotRegistered
	}

	user, err := uc.mapper.RegisterRequestToUser(req)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user entity")
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.userRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user in database")
		return errors.ErrInternalServer
	}

	if err := uc.parentRepo.UpdateUserId(ctx, tx, parent.TempEmail, user.Id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to update parent of user")
		return errors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	verificationToken, err := uc.mapper.CreateVerificationToken(user.Id)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create verification token")
	}

	if verificationToken != nil {
		if err := uc.verifyTokenRepo.Create(ctx, verificationToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
			return errors.ErrInternalServer
		}

		verifyLink := fmt.Sprintf("http://localhost:3000/api/v1/auth/verify-account?token=%s", verificationToken.Token)
		if err := uc.emailService.SendVerificationEmail(user.Email, user.Username, verifyLink); err != nil {
			log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
			return errors.ErrInternalServer
		}
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("User registered and verification email sent")

	return nil
}

func (uc *authUseCase) ResendVerificationAccountUseCase(ctx context.Context, req *dto.ResendTokenRequest) error {
	if err := uc.validator.ValidateResendEmailRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to find user by email")
		return errors.ErrEmailNotFound
	}

	if user.IsActive {
		log.Warn().Str("email", req.Email).Msg("User is already active")
		return errors.ErrEmailAlreadyVerified
	}

	existingToken, err := uc.verifyTokenRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to check existing token")
		return errors.ErrInternalServer
	}

	var verificationToken *entities.VerificationToken
	var verifyLink string

	if existingToken != nil && !existingToken.IsExpired() {
		verificationToken = existingToken
		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/verify-account?token=%s", existingToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Reusing existing valid token")
	} else {
		verificationToken, err = uc.mapper.CreateVerificationToken(user.Id)
		if err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create new verification token")
			return errors.ErrInternalServer
		}

		if err := uc.verifyTokenRepo.Create(ctx, verificationToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save new verification token")
			return errors.ErrInternalServer
		}

		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/verify-account?token=%s", verificationToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Created and saved new verification token")
	}

	if err := uc.emailService.SendVerificationEmail(user.Email, user.Username, verifyLink); err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
		return errors.ErrInternalServer
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Verification email resent successfully")
	return nil
}

func (uc *authUseCase) VerificationAccountUseCase(ctx context.Context, req *dto.VerifyTokenRequest) error {
	if err := uc.validator.VerificationAccountRequest(req); err != nil {
		log.Warn().Err(err).Str("token", req.Token).Msg("Verification account failed")
		return err
	}

	verificationToken, err := uc.verifyTokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		log.Warn().Err(err).Str("code", req.Token).Msg("Invalid verification code")
		return errors.ErrInvalidToken
	}

	if err := uc.verifyTokenRepo.UpdateStatus(ctx, verificationToken.Token); err != nil {
		log.Error().Err(err).Str("code", req.Token).Msg("Failed to update verification code")
		return errors.ErrInternalServer
	}

	if err := uc.parentRepo.UpdateRegistrationStatus(ctx, verificationToken.UserId); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to complete registration")
		return errors.ErrInternalServer
	}

	if err := uc.userRepo.UpdateActiveStatus(ctx, verificationToken.UserId, true); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to activate user")
		return errors.ErrInternalServer
	}

	log.Info().
		Str("userId", verificationToken.UserId).
		Msg("Email verified successfully")

	return nil
}

func (uc *authUseCase) ForgetPasswordUseCase(ctx context.Context, req *dto.ForgetPasswordRequest) error {
	if err := uc.validator.ValidateForgetPasswordRequest(req); err != nil {
		return err
	}

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Reset password attempt failed: email not found")
		return errors.ErrEmailNotFound
	}

	verificationCode, err := uc.mapper.CreateVerificationToken(user.Id)
	if err != nil {
		log.Warn().Err(err).Str("userId", user.Id).Msg("Failed to create forget password code")
		return errors.ErrInternalServer
	}

	if verificationCode != nil {
		if err := uc.verifyTokenRepo.Create(ctx, verificationCode); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
			return errors.ErrInternalServer
		}

		verifyLink := fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", verificationCode.Token)
		if err := uc.emailService.SendResetPasswordEmail(user.Email, user.Username, verifyLink); err != nil {
			log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
			return errors.ErrInternalServer
		}
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Forget password code sent")
	return nil
}

func (uc *authUseCase) ResendForgetPasswordUseCase(ctx context.Context, req *dto.ResendTokenRequest) error {
	if err := uc.validator.ValidateResendEmailRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to find user by email")
		return errors.ErrEmailNotFound
	}

	if user.IsActive {
		log.Warn().Str("email", req.Email).Msg("User is already active")
		return errors.ErrEmailAlreadyVerified
	}

	existingToken, err := uc.verifyTokenRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to check existing token")
		return errors.ErrInternalServer
	}

	var verificationToken *entities.VerificationToken
	var verifyLink string

	if existingToken != nil && !existingToken.IsExpired() {
		verificationToken = existingToken
		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", existingToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Reusing existing valid token")
	} else {
		verificationToken, err = uc.mapper.CreateVerificationToken(user.Id)
		if err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create new verification token")
			return errors.ErrInternalServer
		}

		if err := uc.verifyTokenRepo.Create(ctx, verificationToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save new verification token")
			return errors.ErrInternalServer
		}

		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", verificationToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Created and saved new verification token")
	}

	if err := uc.emailService.SendResetPasswordEmail(user.Email, user.Username, verifyLink); err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
		return errors.ErrInternalServer
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Verification email resent successfully")
	return nil
}

func (uc *authUseCase) ResetPasswordUseCase(ctx context.Context, req *dto.ResetPasswordRequest) error {
	if err := uc.validator.ValidateResetPasswordRequest(req); err != nil {
		log.Warn().Err(err).Msg("Reset password validation failed")
		return err
	}

	token := req.Token

	if token == "" {
		log.Warn().Msg("Verification token is empty")
		return errors.ErrInvalidToken
	}

	verificationToken, err := uc.verifyTokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		log.Warn().Err(err).Str("code", req.Token).Msg("Invalid verification code")
		return errors.ErrInvalidToken
	}

	if verificationToken.IsExpired() {
		log.Warn().Str("code", req.Token).Str("userId", verificationToken.UserId).Msg("Verification code expired")
		return errors.ErrTokenExpired
	}

	user, err := uc.mapper.ResetPasswordRequestToUser(req)
	if err != nil {
		log.Warn().Err(err).Str("userId", user.Id).Msg("Failed to create forget password code")
		return errors.ErrInternalServer
	}

	if req.Password != req.ConfirmPassword {
		return errors.ErrPasswordNotSame
	}

	if err := uc.userRepo.UpdatePassword(ctx, verificationToken.UserId, user.Password); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to update password")
		return errors.ErrInternalServer
	}

	log.Info().
		Str("userId", verificationToken.UserId).
		Msg("Update password successfully")

	return nil
}

func (uc *authUseCase) LoginUseCase(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := uc.validator.ValidateLoginRequest(req); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login validation failed")
		return nil, err
	}

	if err := uc.rateLimiter.CheckLoginRateLimit(ctx, req.Identifier); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login blocked due to rate limiting")
		return nil, errors.ErrTooManyLoginAttempts
	}

	user, err := uc.userRepo.GetByIdentifier(ctx, req.Identifier)
	if err != nil {
		uc.rateLimiter.IncrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("identifier", req.Identifier).Msg("Login attempt failed: user not found")
		return nil, errors.ErrInvalidCredentials
	}

	if !user.IsActive {
		log.Warn().Str("identifier", req.Identifier).Str("userId", user.Id).Msg("Login failed: user inactive")
		return nil, errors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		uc.rateLimiter.IncrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("username", req.Identifier).Str("userId", user.Id).Msg("Login attempt failed: invalid password")
		return nil, errors.ErrInvalidCredentials
	}

	uc.rateLimiter.ClearFailedAttempts(ctx, req.Identifier)

	accessToken, err := uc.tokenService.GenerateAccessToken(user.Id, user.Role)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to generate access token")
		return nil, errors.ErrGenerateToken
	}

	refreshToken, err := uc.mapper.CreateRefreshToken(user.Id)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create refresh token")
	}

	if refreshToken != nil {
		if err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save refresh token")
			return nil, errors.ErrSaveRefreshToken
		}
	}

	response := uc.mapper.LoginResponse(user, refreshToken)
	response.AccessToken = accessToken

	log.Info().
		Str("userId", user.Id).
		Str("username", user.Username).
		Msg("Successfully logged in")

	return response, nil
}

func (uc *authUseCase) RefreshTokenUseCase(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		log.Warn().Msg("Refresh token cannot be nil")
		return nil, errors.ErrInvalidRefreshToken
	}

	refreshToken, err := uc.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		log.Warn().Str("token", req.RefreshToken[:10]+"...").Msg("Invalid refresh token provided")
		return nil, errors.ErrInvalidRefreshToken
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	user, err := uc.userRepo.GetById(ctx, refreshToken.UserId)
	if err != nil {
		log.Error().Err(err).Str("user_id", refreshToken.UserId).Msg("User associated with refresh token not found")
		return nil, errors.ErrUserNotFound
	}

	if !user.IsActive {
		log.Warn().Str("userId", user.Id).Msg("Refresh token request for inactive user")
		return nil, errors.ErrUserInactive
	}

	newAccessToken, err := helpers.GenerateToken(refreshToken.UserId, constants.Role(user.Role))
	if err != nil {
		log.Error().Err(err).Str("user_id", refreshToken.UserId).Msg("Failed to generate new access token during refresh")
		return nil, errors.ErrInternalServer
	}

	response := uc.mapper.RefreshTokenToResponse(refreshToken)
	response.AccessToken = newAccessToken
	response.ExpiresAt = refreshToken.ExpiresAt.Format("2006-01-02 15:04:05")

	log.Info().Str("user_id", refreshToken.UserId).Msg("Access token refreshed successfully")
	return response, nil
}

func (uc *authUseCase) LogoutUseCase(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		log.Warn().Msg("Logout attempt with empty refresh token")
		return errors.ErrInvalidRefreshToken
	}

	if err := uc.refreshTokenRepo.RevokeStatus(ctx, refreshToken); err != nil {
		log.Warn().Str("token", refreshToken[:10]+"...").Msg("Failed to revoke refresh token during logout")
		return errors.ErrInvalidRefreshToken
	}

	log.Info().Str("token", refreshToken[:10]+"...").Msg("User logged out successfully")
	return nil
}
