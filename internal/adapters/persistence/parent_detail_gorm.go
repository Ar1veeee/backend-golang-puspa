package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/pkg/models"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type parentDetailRepository struct {
	db *gorm.DB
}

func NewParentDetailRepository(db *gorm.DB) repositories.ParentDetailRepository {
	return &parentDetailRepository{
		db: db,
	}
}

func (r *parentDetailRepository) Create(ctx context.Context, tx *gorm.DB, parentDetail *entities.ParentDetail) error {
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
