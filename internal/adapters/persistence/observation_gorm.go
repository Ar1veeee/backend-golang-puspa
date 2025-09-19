package persistence

import (
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/database/models"
	"errors"
	"fmt"
	"time"

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

	observations := make([]*entities.Observation, 0, len(dbObservations))
	for _, dbObservation := range dbObservations {
		observation := r.modelToEntity(dbObservation)
		observations = append(observations, observation)
	}

	return observations, nil
}

func (r *observationRepository) GetByScheduledStatus(ctx context.Context) ([]*entities.Observation, error) {
	var dbObservations []*models.Observation

	if err := r.db.WithContext(ctx).
		Preload("Children").
		Preload("Children.Parent").
		Preload("Children.Parent.ParentDetail").
		Where("status = ?", "Scheduled").
		Order("scheduled_date asc").
		Find(&dbObservations).Error; err != nil {
		return nil, fmt.Errorf("failed to get scheduled observations: %w", err)
	}

	observations := make([]*entities.Observation, 0, len(dbObservations))
	for _, dbObservation := range dbObservations {
		observation := r.modelToEntity(dbObservation)
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

	observations := make([]*entities.Observation, 0, len(dbObservations))
	for _, dbObservation := range dbObservations {
		observation := r.modelToEntity(dbObservation)
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

	observation := r.modelToEntity(&dbObservation)
	if observation == nil {
		return nil, errors.New("failed to find observation by id")
	}

	return observation, nil
}

func (r *observationRepository) UpdateScheduledDate(ctx context.Context, observationId int, date helpers.DateOnly) error {
	if observationId == 0 {
		return errors.New("observation is nil")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Observation{}).
		Where("id = ?", observationId).
		Updates(map[string]interface{}{
			"scheduled_date": date,
			"updated_at":     time.Now(),
			"status":         string(constants.ObservationStatusScheduled),
		})

	if result.Error != nil {
		return errors.New("failed to update observation date")
	}

	if result.RowsAffected == 0 {
		return errors.New("observation not found")
	}

	return nil
}

func (r *observationRepository) UpdateAfterObservation(ctx context.Context, tx *gorm.DB, observationId int, therapistId string, totalScore int, conclusion string, recommendation string) error {
	if observationId == 0 {
		return errors.New("observation is nil")
	}

	result := tx.WithContext(ctx).
		Model(&models.Observation{}).
		Where("id = ?", observationId).
		Updates(map[string]interface{}{
			"therapist_id":   therapistId,
			"total_score":    totalScore,
			"conclusion":     conclusion,
			"recommendation": recommendation,
			"status":         string(constants.ObservationStatusCompleted),
			"updated_at":     time.Now(),
		})

	if result.Error != nil {
		return errors.New("failed to update observation")
	}

	if result.RowsAffected == 0 {
		return errors.New("observation not found")
	}

	return nil
}

func (r *observationRepository) modelToEntity(dbObservation *models.Observation) *entities.Observation {
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

func (r *observationRepository) entityToModel(observation *entities.Observation) *models.Observation {
	return &models.Observation{
		Id:             observation.Id,
		ChildId:        observation.ChildId,
		TherapistId:    observation.TherapistId,
		ScheduledDate:  observation.ScheduledDate,
		AgeCategory:    observation.AgeCategory,
		TotalScore:     observation.TotalScore,
		Conclusion:     observation.Conclusion,
		Recommendation: observation.Recommendation,
		Status:         observation.Status,
		CreatedAt:      observation.CreatedAt,
		UpdatedAt:      observation.UpdatedAt,
	}
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
