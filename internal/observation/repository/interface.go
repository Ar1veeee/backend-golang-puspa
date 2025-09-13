package repository

import (
	"backend-golang/internal/observation/entity"
	"context"
)

type ObservationRepository interface {
	GetPendingObservations(ctx context.Context) ([]*entity.Observation, error)
	GetCompletedObservations(ctx context.Context) ([]*entity.Observation, error)
	GetObservationById(ctx context.Context, id int) (*entity.Observation, error)
}
