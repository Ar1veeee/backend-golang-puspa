package gorm

import (
	"backend-golang/internal/observation/entity"
	"backend-golang/internal/observation/repository"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type observationRepository struct {
	db *gorm.DB
}

func NewObservationRepository(db *gorm.DB) repository.ObservationRepository {
	return &observationRepository{db: db}
}

func (r *observationRepository) GetPendingObservations(ctx context.Context) ([]*entity.Observation, error) {
	var dbObservations []*models.Observation

	if err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Children.Parent").
		Preload("Children.Parent.ParentDetail").
		Where("status = ?", "Pending").
		Order("scheduled_date desc").
		Find(&dbObservations).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending observations: %w", err)
	}

	var observations []*entity.Observation
	for _, dbObservation := range dbObservations {
		observation := r.modelToObservationsEntity(dbObservation)
		observations = append(observations, observation)
	}

	return observations, nil
}

func (r *observationRepository) GetCompletedObservations(ctx context.Context) ([]*entity.Observation, error) {
	var dbObservations []*models.Observation

	if err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Children.Parent").
		Preload("Children.Parent.ParentDetail").
		Where("status = ?", "Complete").
		Order("scheduled_date asc").
		Find(&dbObservations).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed observations: %w", err)
	}

	var observations []*entity.Observation
	for _, dbObservation := range dbObservations {
		observation := r.modelToObservationsEntity(dbObservation)
		observations = append(observations, observation)
	}

	return observations, nil
}

func (r *observationRepository) GetObservationById(ctx context.Context, id int) (*entity.Observation, error) {
	if id == 0 {
		return nil, errors.New("id cannot be empty")
	}

	var dbObservations models.Observation

	if err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Children.Parent").
		Preload("Children.Parent.ParentDetail").
		First(&dbObservations, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("observation not found")
		}
		return nil, errors.New("failed to find observation by id")
	}

	observation := r.modelToObservationsEntity(&dbObservations)
	if observation == nil {
		return nil, errors.New("failed to find observation by id")
	}

	return observation, nil
}

func (r *observationRepository) modelToParentEntity(dbParent *models.Parent) *entity.Parent {
	parent := &entity.Parent{
		Id:                 dbParent.Id,
		TempEmail:          dbParent.TempEmail,
		RegistrationStatus: dbParent.RegistrationStatus,
		CreatedAt:          dbParent.CreatedAt,
		UpdatedAt:          dbParent.UpdatedAt,
	}

	if dbParent.UserId != nil {
		parent.UserId = dbParent.UserId
	}

	if len(dbParent.ParentDetail) > 0 {
		var parentDetails []entity.ParentDetail
		for _, parentDetail := range dbParent.ParentDetail {
			if converted := r.modelToParentDetailEntity(&parentDetail); converted != nil {
				parentDetails = append(parentDetails, *converted)
			}
		}
		parent.ParentDetail = parentDetails
	}

	return parent
}

func (r *observationRepository) modelToParentDetailEntity(dbParentDetail *models.ParentDetail) *entity.ParentDetail {
	parentDetail := &entity.ParentDetail{
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

func (r *observationRepository) modelToChildrenEntity(dbChildren *models.Children) *entity.Children {
	child := &entity.Children{
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

func (r *observationRepository) modelToObservationsEntity(dbObservation *models.Observation) *entity.Observation {
	observation := &entity.Observation{
		Id:             dbObservation.Id,
		ChildId:        dbObservation.ChildId,
		TherapistId:    dbObservation.TherapistId,
		ScheduledDate:  dbObservation.ScheduledDate,
		AgeCategory:    dbObservation.AgeCategory,
		TotalScore:     dbObservation.TotalScore,
		Conclusion:     dbObservation.Conclusion,
		Recommendation: dbObservation.Recommendation,
		Status:         dbObservation.Status,
		CreatedAt:      dbObservation.CreatedAt,
		UpdatedAt:      dbObservation.UpdatedAt,
	}

	if dbObservation.Children != nil {
		observation.Children = r.modelToChildrenEntity(dbObservation.Children)
	}

	return observation
}
