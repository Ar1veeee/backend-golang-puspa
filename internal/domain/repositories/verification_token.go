package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"
)

type VerificationTokenRepository interface {
	Create(ctx context.Context, token *entities.VerificationToken) error
	GetByEmail(ctx context.Context, email string) (*entities.VerificationToken, error)
	GetByToken(ctx context.Context, token string) (*entities.VerificationToken, error)
	UpdateStatus(ctx context.Context, token string) error
}
