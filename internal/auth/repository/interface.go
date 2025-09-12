package repository

import (
	"backend-golang/internal/auth/entity"
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByIdentifier(ctx context.Context, identifier string) (*entity.User, error)
	FindUserById(ctx context.Context, id string) (*entity.User, error)
	FindTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateUserActiveStatus(ctx context.Context, userID string, isActive bool) error
	UpdateParentRegistrationStatus(ctx context.Context, userId string) error
	ResetUserPassword(ctx context.Context, userId, newPassword string) error

	SaveVerificationToken(ctx context.Context, code *entity.VerificationToken) error
	VerifyAccountByToken(ctx context.Context, code string) (*entity.VerificationToken, error)

	SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error

	GetParentByTempEmail(ctx context.Context, email string) (*entity.Parent, error)
	UpdateParentUserId(ctx context.Context, tempEmail string, userID string) error
}
