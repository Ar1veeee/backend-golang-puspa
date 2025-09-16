package repositories

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Begin(ctx context.Context) *gorm.DB
}
