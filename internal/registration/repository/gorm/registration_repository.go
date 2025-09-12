package gorm

import (
	"backend-golang/internal/registration/entity"
	"backend-golang/internal/registration/repository"
	"backend-golang/shared/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) repository.RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *registrationRepository) CreateParentWithTx(ctx context.Context, tx *gorm.DB, parent *entity.Parent) error {
	dbParent := &models.Parent{
		Id:                 parent.Id,
		UserId:             nil,
		TempEmail:          parent.TempEmail,
		RegistrationStatus: parent.RegistrationStatus,
		CreatedAt:          parent.CreatedAt,
		UpdatedAt:          parent.UpdatedAt,
	}
	if err := tx.WithContext(ctx).Create(&dbParent).Error; err != nil {
		return fmt.Errorf("failed to create parent: %w", err)
	}

	return nil
}

func (r *registrationRepository) CreateParentDetailWithTx(ctx context.Context, tx *gorm.DB, parentDetail *entity.ParentDetail) error {
	dbParentDetail := &models.ParentDetail{
		Id:                    parentDetail.Id,
		ParentId:              parentDetail.ParentId,
		ParentType:            parentDetail.ParentType,
		ParentName:            parentDetail.ParentName,
		ParentPhone:           parentDetail.ParentPhone,
		ParentBirthDate:       nil,
		ParentOccupation:      nil,
		RelationshipWithChild: nil,
	}

	if err := tx.WithContext(ctx).Create(&dbParentDetail).Error; err != nil {
		return fmt.Errorf("failed to create parent detail: %w", err)
	}

	return nil
}

func (r *registrationRepository) CreateChildWithTx(ctx context.Context, tx *gorm.DB, child *entity.Children) error {
	dbChild := &models.Children{
		Id:                 child.Id,
		ParentId:           child.ParentId,
		ChildName:          child.ChildName,
		ChildGender:        child.ChildGender,
		ChildBirthPlace:    child.ChildBirthPlace,
		ChildBirthDate:     child.ChildBirthDate,
		ChildAddress:       child.ChildAddress,
		ChildComplaint:     child.ChildComplaint,
		ChildSchool:        child.ChildSchool,
		ChildServiceChoice: child.ChildServiceChoice,
		ChildReligion:      nil,
		CreatedAt:          child.CreatedAt,
		UpdatedAt:          child.UpdatedAt,
	}
	if err := tx.WithContext(ctx).Create(&dbChild).Error; err != nil {
		return fmt.Errorf("failed to create children: %w", err)
	}

	return nil
}

func (r *registrationRepository) CreateObservationWithTx(ctx context.Context, tx *gorm.DB, observation *entity.Observation) error {
	dbObservation := &models.Observation{
		Id:            observation.Id,
		ChildId:       observation.ChildId,
		ScheduledDate: observation.ScheduledDate,
		AgeCategory:   observation.AgeCategory,
		Status:        observation.Status,
		CreatedAt:     observation.CreatedAt,
		UpdatedAt:     observation.UpdatedAt,
	}
	if err := tx.WithContext(ctx).Create(&dbObservation).Error; err != nil {
		return fmt.Errorf("failed to create observation: %w", err)
	}

	return nil
}

func (r *registrationRepository) ExistsByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.Parent{}).Where("temp_email = ?", email).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
