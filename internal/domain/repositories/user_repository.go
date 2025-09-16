package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByIdentifier(ctx context.Context, identifier string) (*entities.User, error)
	GetById(ctx context.Context, id string) (*entities.User, error)

	Update(ctx context.Context, tx *gorm.DB, user *entities.User) error
	UpdateActiveStatus(ctx context.Context, userID string, isActive bool) error
	UpdatePassword(ctx context.Context, userId, newPassword string) error

	CheckExisting(ctx context.Context, email, username string) (emailExists, usernameExists bool, err error)
}
