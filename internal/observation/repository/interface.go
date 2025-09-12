package repository

import (
	"backend-golang/internal/observation/entity"
	"context"
)

type ObservationRepository interface {
	GetPendingObservations(ctx context.Context) ([]*entity.Observation, error)
}
