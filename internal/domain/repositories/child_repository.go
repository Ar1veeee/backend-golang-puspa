package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ChildRepository interface {
	Create(ctx context.Context, tx *gorm.DB, child *entities.Children) error
}
