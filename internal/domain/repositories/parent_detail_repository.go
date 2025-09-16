package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ParentDetailRepository interface {
	Create(ctx context.Context, tx *gorm.DB, child *entities.ParentDetail) error
}
