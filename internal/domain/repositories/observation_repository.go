package repositories

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/helpers"
	"context"

	"gorm.io/gorm"
)

type ObservationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, child *entities.Observation) error

	GetByPendingStatus(ctx context.Context) ([]*entities.Observation, error)
	GetByScheduledStatus(ctx context.Context) ([]*entities.Observation, error)
	GetByCompletedStatus(ctx context.Context) ([]*entities.Observation, error)
	GetById(ctx context.Context, observationId int) (*entities.Observation, error)

	UpdateScheduledDate(ctx context.Context, observationId int, date helpers.DateOnly) error
	UpdateAfterObservation(ctx context.Context, tx *gorm.DB, observationId int, therapistId string, totalScore int, conclusion string, recommendation string) error
}
