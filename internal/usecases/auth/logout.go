package auth

import (
	"backend-golang/internal/errors"
	"context"

	"github.com/rs/zerolog/log"
)

type logoutUseCase struct {
	deps *Dependencies
}

func NewLogoutUseCase(deps *Dependencies) LogoutUseCase {
	return &logoutUseCase{deps: deps}
}

func (uc *logoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		log.Warn().Msg("Logout attempt with empty refresh token")
		return errors.ErrInvalidRefreshToken
	}

	if err := uc.deps.RefreshTokenRepo.RevokeStatus(ctx, refreshToken); err != nil {
		log.Warn().Str("token", refreshToken[:10]+"...").Msg("Failed to revoke refresh token during logout")
		return errors.ErrInvalidRefreshToken
	}

	log.Info().Str("token", refreshToken[:10]+"...").Msg("User logged out successfully")
	return nil
}
