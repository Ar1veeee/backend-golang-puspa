package admin

import (
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type deleteAdminUseCase struct {
	deps *Dependencies
}

func NewDeleteAdminUseCase(deps *Dependencies) DeleteAdminUseCase {
	return &deleteAdminUseCase{deps: deps}
}

func (uc *deleteAdminUseCase) Execute(ctx context.Context, adminId string) error {
	if adminId == "" {
		return fmt.Errorf("adminId is empty")
	}

	admin, err := uc.deps.AdminRepo.GetById(ctx, adminId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	if err := uc.deps.AdminRepo.Delete(ctx, tx, admin); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDeletionFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
