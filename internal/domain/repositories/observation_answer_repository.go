package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ObservationAnswerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, answers []*entities.ObservationAnswer) error
}
