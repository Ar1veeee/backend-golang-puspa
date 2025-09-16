package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(ctx context.Context, tx *gorm.DB, therapist *entities.Admin) error

	GetAll(ctx context.Context) ([]*entities.Admin, error)
	GetById(ctx context.Context, adminId string) (*entities.Admin, error)

	Update(ctx context.Context, tx *gorm.DB, admin *entities.Admin) error

	Delete(ctx context.Context, tx *gorm.DB, admin *entities.Admin) error
}
