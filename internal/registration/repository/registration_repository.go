package repository

import (
	"backend-golang/internal/registration/entity"
	"context"

	"gorm.io/gorm"
)

type RegistrationRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	CreateParentWithTx(ctx context.Context, tx *gorm.DB, parent *entity.Parent) error
	CreateParentDetailWithTx(ctx context.Context, tx *gorm.DB, parentDetail *entity.ParentDetail) error
	CreateChildWithTx(ctx context.Context, tx *gorm.DB, child *entity.Children) error
	CreateObservationWithTx(ctx context.Context, tx *gorm.DB, observation *entity.Observation) error
	ExistsByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)
}
