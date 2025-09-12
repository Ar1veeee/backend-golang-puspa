package gorm

import (
	"backend-golang/internal/therapist/entity"
	"backend-golang/internal/therapist/repository"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type therapistRepository struct {
	db *gorm.DB
}

func NewTherapistRepository(db *gorm.DB) repository.TherapistRepository {
	return &therapistRepository{db: db}
}

func (r *therapistRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *therapistRepository) CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if user == nil {
		return errors.New("user data cannot be empty")
	}

	dbUser := &models.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		IsActive:  true,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := tx.WithContext(ctx).Create(&dbUser).Error; err != nil {
		return fmt.Errorf("failed create user: %w", err)
	}

	return nil
}

func (r *therapistRepository) CreateTherapistWithTx(ctx context.Context, tx *gorm.DB, therapist *entity.Therapist) error {
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

func (r *therapistRepository) GetAllTherapist(ctx context.Context) ([]*entity.Therapist, error) {
	var dbTherapists []*models.Therapist

	if err := r.db.WithContext(ctx).
		Preload("User").
		Find(&dbTherapists).Error; err != nil {
		return nil, fmt.Errorf("failed to get therapists: %w", err)
	}

	var therapists []*entity.Therapist
	for _, dbTherapist := range dbTherapists {
		therapist := r.modelToTherapistEntity(dbTherapist)
		therapists = append(therapists, therapist)
	}

	return therapists, nil
}

func (r *therapistRepository) GetByTherapistId(ctx context.Context, id string) (*entity.Therapist, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var dbTherapist models.Therapist

	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&dbTherapist, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("therapist with id %s not found", id)
		}
		return nil, fmt.Errorf("failed get therapist with id %s: %w", id, err)
	}

	return r.modelToTherapistEntity(&dbTherapist), nil
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

func (r *therapistRepository) modelToUserEntity(dbUser *models.User) *entity.User {
	return &entity.User{
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

func (r *therapistRepository) modelToTherapistEntity(dbTherapist *models.Therapist) *entity.Therapist {
	therapist := &entity.Therapist{
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
	}

	return therapist
}
