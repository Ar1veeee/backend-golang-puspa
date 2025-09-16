package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ObservationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, child *entities.Observation) error
	GetByPendingStatus(ctx context.Context) ([]*entities.Observation, error)
	GetByCompletedStatus(ctx context.Context) ([]*entities.Observation, error)
	GetById(ctx context.Context, observationId int) (*entities.Observation, error)
}
