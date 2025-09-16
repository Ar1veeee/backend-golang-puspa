package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"context"

	"backend-golang/pkg/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type therapistRepository struct {
	db *gorm.DB
}

func NewTherapistRepository(db *gorm.DB) repositories.TherapistRepository {
	return &therapistRepository{
		db: db,
	}
}

func (r *therapistRepository) Create(ctx context.Context, tx *gorm.DB, therapist *entities.Therapist) error {
	if therapist == nil {
		return errors.New("therapist data cannot be empty")
	}

	var dbTherapist = &models.Therapist{
		Id:               therapist.Id,
		UserId:           therapist.UserId,
		TherapistName:    therapist.TherapistName,
		TherapistSection: therapist.TherapistSection,
		TherapistPhone:   therapist.TherapistPhone,
		CreatedAt:        therapist.CreatedAt,
		UpdatedAt:        therapist.UpdatedAt,
	}

	if err := tx.WithContext(ctx).Create(&dbTherapist).Error; err != nil {
		return fmt.Errorf("failed create therapist: %w", err)
	}

	return nil
}

func (r *therapistRepository) GetAll(ctx context.Context) ([]*entities.Therapist, error) {
	var dbTherapists []*models.Therapist

	if err := r.db.WithContext(ctx).
		Preload("User", "is_active = ?", true).
		Find(&dbTherapists).Error; err != nil {
		return nil, fmt.Errorf("failed to get therapists: %w", err)
	}

	therapists := make([]*entities.Therapist, 0, len(dbTherapists))
	for _, dbTherapist := range dbTherapists {
		therapist := r.modelToTherapistEntity(dbTherapist)
		therapists = append(therapists, therapist)
	}

	return therapists, nil
}

func (r *therapistRepository) GetById(ctx context.Context, therapistId string) (*entities.Therapist, error) {
	var dbTherapist models.Therapist

	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", therapistId).
		First(&dbTherapist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("therapist with id %s not found", therapistId)
		}
		return nil, fmt.Errorf("failed get therapist with id %s: %w", therapistId, err)
	}

	therapist := r.modelToTherapistEntity(&dbTherapist)
	if therapist == nil {
		return nil, errors.New("failed to find therapist by id")
	}

	return therapist, nil
}

func (r *therapistRepository) Update(ctx context.Context, tx *gorm.DB, therapist *entities.Therapist) error {
	if therapist == nil {
		return errors.New("therapist data cannot be empty")
	}

	dbTherapist := &models.Therapist{
		Id:               therapist.Id,
		UserId:           therapist.UserId,
		TherapistName:    therapist.TherapistName,
		TherapistSection: therapist.TherapistSection,
		TherapistPhone:   therapist.TherapistPhone,
		CreatedAt:        therapist.CreatedAt,
		UpdatedAt:        therapist.UpdatedAt,
	}

	if err := tx.WithContext(ctx).Save(dbTherapist).Error; err != nil {
		return fmt.Errorf("failed to update therapist: %w", err)
	}

	return nil
}

func (r *therapistRepository) Delete(ctx context.Context, tx *gorm.DB, therapist *entities.Therapist) error {
	if therapist == nil {
		return errors.New("therapist data cannot be empty")
	}

	if therapist.User != nil {
		if err := tx.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", therapist.User.Id).
			Update("is_active", false).Error; err != nil {
			return fmt.Errorf("failed to deactive therapist: %w", err)
		}
	}

	if err := tx.WithContext(ctx).
		Model(&models.Therapist{}).
		Where("id = ?", therapist.Id).
		Update("is_deleted", true).Error; err != nil {
		return fmt.Errorf("failed to delete therapist: %w", err)
	}

	return nil
}

func (r *therapistRepository) modelToUserEntity(dbUser *models.User) *entities.User {
	return &entities.User{
		Id:        dbUser.Id,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func (r *therapistRepository) modelToTherapistEntity(dbTherapist *models.Therapist) *entities.Therapist {
	therapist := &entities.Therapist{
		Id:               dbTherapist.Id,
		UserId:           dbTherapist.UserId,
		TherapistName:    dbTherapist.TherapistName,
		TherapistSection: dbTherapist.TherapistSection,
		TherapistPhone:   dbTherapist.TherapistPhone,
		CreatedAt:        dbTherapist.CreatedAt,
		UpdatedAt:        dbTherapist.UpdatedAt,
	}

	if dbTherapist.User != nil {
		therapist.User = r.modelToUserEntity(dbTherapist.User)
	} else {
		therapist.User = nil
	}

	return therapist
}
