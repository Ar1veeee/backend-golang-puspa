package persistence

import (
	"backend-golang/internal/domain/repositories"
	"context"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) Begin(ctx context.Context) *gorm.DB {
	return t.db.WithContext(ctx).Begin()
}
