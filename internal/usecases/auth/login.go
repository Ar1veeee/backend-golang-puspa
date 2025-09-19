package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type loginUseCase struct {
	deps *Dependencies
}

func NewLoginUseCase(deps *Dependencies) LoginUseCase {
	return &loginUseCase{deps: deps}
}

func (uc *loginUseCase) Execute(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := uc.deps.Validator.ValidateLoginRequest(req); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login validation failed")
		return nil, err
	}

	if err := uc.deps.RateLimiter.CheckLoginRateLimit(ctx, req.Identifier); err != nil {
		log.Warn().Err(err).Str("identifier", req.Identifier).Msg("Login blocked due to rate limiting")
		return nil, errors.ErrTooManyLoginAttempts
	}

	user, err := uc.deps.UserRepo.GetByIdentifier(ctx, req.Identifier)
	if err != nil {
		uc.deps.RateLimiter.IncrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("identifier", req.Identifier).Msg("Login attempt failed: user not found")
		return nil, errors.ErrInvalidCredentials
	}

	if !user.IsActive {
		log.Warn().Str("identifier", req.Identifier).Str("userId", user.Id).Msg("Login failed: user inactive")
		return nil, errors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		uc.deps.RateLimiter.IncrementFailedAttempts(ctx, req.Identifier)
		log.Warn().Str("username", req.Identifier).Str("userId", user.Id).Msg("Login attempt failed: invalid password")
		return nil, errors.ErrInvalidCredentials
	}

	uc.deps.RateLimiter.ClearFailedAttempts(ctx, req.Identifier)

	accessToken, err := uc.deps.TokenService.GenerateAccessToken(user.Id, user.Role)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to generate access token")
		return nil, errors.ErrGenerateToken
	}

	refreshToken, err := uc.deps.Mapper.CreateRefreshToken(user.Id)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create refresh token")
	}

	if refreshToken != nil {
		if err := uc.deps.RefreshTokenRepo.Create(ctx, refreshToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save refresh token")
			return nil, errors.ErrSaveRefreshToken
		}
	}

	response := uc.deps.Mapper.LoginResponse(user, refreshToken)
	response.AccessToken = accessToken

	log.Info().
		Str("userId", user.Id).
		Str("username", user.Username).
		Msg("Successfully logged in")

	return response, nil
}
