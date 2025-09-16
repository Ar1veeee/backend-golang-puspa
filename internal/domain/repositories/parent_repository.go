package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ParentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, parent *entities.Parent) error
	GetByTempEmail(ctx context.Context, email string) (*entities.Parent, error)

	UpdateRegistrationStatus(ctx context.Context, userId string) error
	UpdateUserId(ctx context.Context, tx *gorm.DB, tempEmail string, userID string) error

	ExistByTempEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)
}
