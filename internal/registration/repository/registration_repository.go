package repository

import (
	"backend-golang/shared/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type registrationRepository struct {
	db *gorm.DB
}

type RegistrationRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	CreateParentWithTx(ctx context.Context, tx *gorm.DB, parent *models.Parent) error
	CreateParentDetailWithTx(ctx context.Context, tx *gorm.DB, parentDetail *models.ParentDetail) error
	CreateChildWithTx(ctx context.Context, tx *gorm.DB, child *models.Children) error
	CreateObservationWithTx(ctx context.Context, tx *gorm.DB, observation *models.Observation) error
	ExistsByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.db.Begin()
}

func (r *registrationRepository) CreateParentWithTx(ctx context.Context, tx *gorm.DB, parent *models.Parent) error {
	return tx.Create(parent).Error
}

func (r *registrationRepository) CreateParentDetailWithTx(ctx context.Context, tx *gorm.DB, parentDetail *models.ParentDetail) error {
	return tx.Create(parentDetail).Error
}

func (r *registrationRepository) CreateChildWithTx(ctx context.Context, tx *gorm.DB, child *models.Children) error {
	return tx.Create(child).Error
}

func (r *registrationRepository) CreateObservationWithTx(ctx context.Context, tx *gorm.DB, observation *models.Observation) error {
	return tx.Create(observation).Error
}

func (r *registrationRepository) ExistsByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.Parent{}).Where("temp_email = ?", email).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
