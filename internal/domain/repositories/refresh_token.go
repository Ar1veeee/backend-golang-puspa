package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	RevokeStatus(ctx context.Context, token string) error
}
