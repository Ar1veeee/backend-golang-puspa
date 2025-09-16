package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/pkg/models"
	"errors"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type observationRepository struct {
	db *gorm.DB
}

func NewObservationRepository(db *gorm.DB) repositories.ObservationRepository {
	return &observationRepository{
		db: db,
	}
}

func (r *observationRepository) Create(ctx context.Context, tx *gorm.DB, observation *entities.Observation) error {
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

func (r *observationRepository) GetByPendingStatus(ctx context.Context) ([]*entities.Observation, error) {
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

	var observations []*entities.Observation
	for _, dbObservation := range dbObservations {
		observation := r.modelToObservationsEntity(dbObservation)
		observations = append(observations, observation)
	}

	return observations, nil
}

func (r *observationRepository) GetByCompletedStatus(ctx context.Context) ([]*entities.Observation, error) {
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

	var observations []*entities.Observation
	for _, dbObservation := range dbObservations {
		observation := r.modelToObservationsEntity(dbObservation)
		observations = append(observations, observation)
	}

	return observations, nil
}

func (r *observationRepository) GetById(ctx context.Context, observationId int) (*entities.Observation, error) {
	if observationId == 0 {
		return nil, errors.New("observationId cannot be empty")
	}

	var dbObservation models.Observation

	if err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Children.Parent").
		Preload("Children.Parent.ParentDetail").
		First(&dbObservation, "id = ?", observationId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("observation not found")
		}
		return nil, errors.New("failed to find observation by id")
	}

	observation := r.modelToObservationsEntity(&dbObservation)
	if observation == nil {
		return nil, errors.New("failed to find observation by id")
	}

	return observation, nil
}

func (r *observationRepository) modelToParentEntity(dbParent *models.Parent) *entities.Parent {
	parent := &entities.Parent{
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

func (r *observationRepository) modelToParentDetailEntity(dbParentDetail *models.ParentDetail) *entities.ParentDetail {
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

func (r *observationRepository) modelToChildrenEntity(dbChildren *models.Children) *entities.Children {
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

func (r *observationRepository) modelToObservationsEntity(dbObservation *models.Observation) *entities.Observation {
	observation := &entities.Observation{
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
