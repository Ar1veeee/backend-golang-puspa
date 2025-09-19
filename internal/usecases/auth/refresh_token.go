package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"

	"github.com/rs/zerolog/log"
)

type refreshTokenUseCase struct {
	deps *Dependencies
}

func NewRefreshTokenUseCase(deps *Dependencies) RefreshTokenUseCase {
	return &refreshTokenUseCase{deps: deps}
}

func (uc *refreshTokenUseCase) Execute(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		log.Warn().Msg("Refresh token cannot be nil")
		return nil, errors.ErrInvalidRefreshToken
	}

	refreshToken, err := uc.deps.RefreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		log.Warn().Str("token", req.RefreshToken[:10]+"...").Msg("Invalid refresh token provided")
		return nil, errors.ErrInvalidRefreshToken
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	user, err := uc.deps.UserRepo.GetById(ctx, refreshToken.UserId)
	if err != nil {
		log.Error().Err(err).Str("user_id", refreshToken.UserId).Msg("User associated with refresh token not found")
		return nil, errors.ErrUserNotFound
	}

	if !user.IsActive {
		log.Warn().Str("userId", user.Id).Msg("Refresh token request for inactive user")
		return nil, errors.ErrUserInactive
	}

	newAccessToken, err := uc.deps.TokenService.GenerateAccessToken(user.Id, user.Role)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to generate access token")
		return nil, errors.ErrGenerateToken
	}

	response := uc.deps.Mapper.RefreshTokenToResponse(refreshToken)
	response.AccessToken = newAccessToken
	response.ExpiresAt = refreshToken.ExpiresAt.Format("2006-01-02 15:04:05")

	log.Info().Str("user_id", refreshToken.UserId).Msg("Access token refreshed successfully")
	return response, nil
}
