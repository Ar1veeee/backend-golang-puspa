package repository

import (
	"backend-golang/internal/auth/entity"
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByUsernameAndEmail(ctx context.Context, identifier string) (*entity.User, error)
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateUserActiveStatus(ctx context.Context, userID string, isActive bool) error

	SaveVerificationCode(ctx context.Context, code *entity.VerificationCode) error
	VerifyEmailByCode(ctx context.Context, code string) (*entity.VerificationCode, error)

	SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error

	GetParentByTempEmail(ctx context.Context, email string) (*entity.Parent, error)
	UpdateParentUserId(ctx context.Context, tempEmail string, userID string) error
}
