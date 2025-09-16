package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/pkg/models"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type childRepository struct {
	db *gorm.DB
}

func NewChildRepository(db *gorm.DB) repositories.ChildRepository {
	return &childRepository{
		db: db,
	}
}

func (r *childRepository) Create(ctx context.Context, tx *gorm.DB, child *entities.Children) error {
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
