package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	models2 "backend-golang/internal/infrastructure/database/models"
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
	dbChild := &models2.Children{
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

//func (r *childRepository) GetById(ctx context.Context, childId string) (*entities.Children, error) {
//	var dbChild *models.Children
//
//	if err := r.db.WithContext(ctx).
//		Preload("Parent").
//		Preload("ParentDetail").
//		Where("child_id = ?", childId).
//		First(&dbChild).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, fmt.Errorf("child not found")
//		}
//		return nil, fmt.Errorf("failed to get children by id: %w", err)
//	}
//
//	child := r.modelToEntity(dbChild)
//	if child == nil {
//		return nil, fmt.Errorf("failed to get children by id")
//	}
//
//	return child, nil
//}

func (r *childRepository) GetAll(ctx context.Context) ([]*entities.Children, error) {
	var dbChilds []*models2.Children

	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Parent.ParentDetail").
		Order("created_at desc").
		Find(&dbChilds).Error; err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}

	children := make([]*entities.Children, 0, len(dbChilds))
	for _, dbChild := range dbChilds {
		child := r.modelToEntity(dbChild)
		children = append(children, child)
	}

	return children, nil
}

func (r *childRepository) modelToEntity(dbChildren *models2.Children) *entities.Children {
	child := &entities.Children{
		Id:                 dbChildren.Id,
		ParentId:           dbChildren.ParentId,
		ChildName:          dbChildren.ChildName,
		ChildGender:        dbChildren.ChildGender,
		ChildBirthPlace:    dbChildren.ChildBirthPlace,
		ChildBirthDate:     dbChildren.ChildBirthDate,
		ChildAddress:       dbChildren.ChildAddress,
		ChildComplaint:     dbChildren.ChildComplaint,
		ChildSchool:        dbChildren.ChildSchool,
		ChildServiceChoice: dbChildren.ChildServiceChoice,
		CreatedAt:          dbChildren.CreatedAt,
		UpdatedAt:          dbChildren.UpdatedAt,
	}

	if dbChildren.Parent != nil {
		child.Parent = r.modelToParentEntity(dbChildren.Parent)
	}

	return child
}

func (r *childRepository) modelToParentEntity(dbParent *models2.Parent) *entities.Parent {
	parent := &entities.Parent{
		Id:                 dbParent.Id,
		TempEmail:          dbParent.TempEmail,
		RegistrationStatus: dbParent.RegistrationStatus,
		CreatedAt:          dbParent.CreatedAt,
		UpdatedAt:          dbParent.UpdatedAt,
	}

	if len(dbParent.ParentDetail) > 0 {
		var parentDetails []entities.ParentDetail
		for _, parentDetail := range dbParent.ParentDetail {
			if converted := r.modelToParentDetailEntity(&parentDetail); converted != nil {
				parentDetails = append(parentDetails, *converted)
			}
		}
		parent.ParentDetail = parentDetails
	}

	return parent
}

func (r *childRepository) modelToParentDetailEntity(dbParentDetail *models2.ParentDetail) *entities.ParentDetail {
	parentDetail := &entities.ParentDetail{
		Id:          dbParentDetail.Id,
		ParentId:    dbParentDetail.ParentId,
		ParentType:  dbParentDetail.ParentType,
		ParentName:  dbParentDetail.ParentName,
		ParentPhone: dbParentDetail.ParentPhone,
		CreatedAt:   dbParentDetail.CreatedAt,
		UpdatedAt:   dbParentDetail.UpdatedAt,
	}

	return parentDetail
}
