package repository

import (
	"backend-golang/internal/therapist/entity"
	"context"

	"gorm.io/gorm"
)

type TherapistRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *entity.User) error
	CreateTherapistWithTx(ctx context.Context, tx *gorm.DB, therapist *entity.Therapist) error
	GetAllTherapist(ctx context.Context) ([]*entity.Therapist, error)
	GetByTherapistId(ctx context.Context, id string) (*entity.Therapist, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
