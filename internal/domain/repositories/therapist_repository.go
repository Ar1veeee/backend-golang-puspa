package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type TherapistRepository interface {
	Create(ctx context.Context, tx *gorm.DB, therapist *entities.Therapist) error

	GetAll(ctx context.Context) ([]*entities.Therapist, error)
	GetById(ctx context.Context, adminId string) (*entities.Therapist, error)

	Update(ctx context.Context, tx *gorm.DB, admin *entities.Therapist) error

	Delete(ctx context.Context, tx *gorm.DB, admin *entities.Therapist) error
}
