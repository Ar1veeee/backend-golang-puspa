package persistence

import (
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"context"

	"backend-golang/pkg/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type parentRepository struct {
	db *gorm.DB
}

func NewParentRepository(db *gorm.DB) repositories.ParentRepository {
	return &parentRepository{db: db}
}

func (r *parentRepository) Create(ctx context.Context, tx *gorm.DB, parent *entities.Parent) error {
	dbParent := r.domainToModel(parent)

	if err := tx.WithContext(ctx).Create(&dbParent).Error; err != nil {
		return fmt.Errorf("failed to create parent: %w", err)
	}

	return nil
}

func (r *parentRepository) UpdateRegistrationStatus(ctx context.Context, userId string) error {
	if userId == "" {
		return errors.New("user id cannot be empty")
	}

	result := r.db.WithContext(ctx).Model(&models.Parent{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"registration_status": string(constants.RegistrationStatusComplete),
			"temp_email":          "",
		})

	if result.Error != nil {
		return errors.New("failed to update registration status")
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *parentRepository) GetByTempEmail(ctx context.Context, email string) (*entities.Parent, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var dbParent models.Parent
	if err := r.db.WithContext(ctx).Where("temp_email = ?", email).First(&dbParent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("failed to find parent by temp email")
	}

	return r.modelToParentDomain(&dbParent), nil
}

func (r *parentRepository) UpdateUserId(ctx context.Context, tx *gorm.DB, tempEmail string, userId string) error {
	if tempEmail == "" || userId == "" {
		return errors.New("user_id or temp_email cannot be empty")
	}

	result := tx.WithContext(ctx).Model(&models.Parent{}).
		Where("temp_email = ?", tempEmail).
		Update("user_id", userId)

	if result.Error != nil {
		return fmt.Errorf("failed to update parent user_id: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("parent with temp_email %s not found", tempEmail)
	}

	return nil
}

func (r *parentRepository) ExistByTempEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.Parent{}).Where("temp_email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

func (r *parentRepository) modelToParentDomain(dbParent *models.Parent) *entities.Parent {
	return &entities.Parent{
		Id:                 dbParent.Id,
		UserId:             dbParent.UserId,
		TempEmail:          dbParent.TempEmail,
		RegistrationStatus: dbParent.RegistrationStatus,
		CreatedAt:          dbParent.CreatedAt,
		UpdatedAt:          dbParent.UpdatedAt,
	}
}

func (r *parentRepository) domainToModel(parent *entities.Parent) *models.Parent {
	return &models.Parent{
		Id:                 parent.Id,
		UserId:             nil,
		TempEmail:          parent.TempEmail,
		RegistrationStatus: parent.RegistrationStatus,
		CreatedAt:          parent.CreatedAt,
		UpdatedAt:          parent.UpdatedAt,
	}
}
