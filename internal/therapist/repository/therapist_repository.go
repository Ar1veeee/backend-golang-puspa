package repository

import (
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type therapistRepository struct {
	db *gorm.DB
}

type TherapistRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error
	CreateTherapistWithTx(ctx context.Context, tx *gorm.DB, therapist *models.Therapist) error
	GetAll(ctx context.Context) ([]*models.Therapist, error)
	GetById(ctx context.Context, id string) (*models.Therapist, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

func NewTherapistRepository(db *gorm.DB) TherapistRepository {
	return &therapistRepository{db: db}
}

func (r *therapistRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.db.Begin()
}

func (r *therapistRepository) CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (r *therapistRepository) CreateTherapistWithTx(ctx context.Context, tx *gorm.DB, therapist *models.Therapist) error {
	return tx.Create(therapist).Error
}

func (r *therapistRepository) GetAll(ctx context.Context) ([]*models.Therapist, error) {
	var therapists []*models.Therapist

	if err := r.db.WithContext(ctx).
		Preload("User").
		Find(&therapists).Error; err != nil {
		return nil, fmt.Errorf("failed to get therapists: %w", err)
	}

	return therapists, nil
}

func (r *therapistRepository) GetById(ctx context.Context, id string) (*models.Therapist, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var therapist models.Therapist

	err := r.db.WithContext(ctx).First(&therapist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &therapist, nil
}

func (r *therapistRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *therapistRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}
