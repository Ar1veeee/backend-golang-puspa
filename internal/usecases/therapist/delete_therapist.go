package therapist

import (
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type deleteTherapistUseCase struct {
	deps *Dependencies
}

func NewDeleteTherapistUseCase(deps *Dependencies) DeleteTherapistUseCase {
	return &deleteTherapistUseCase{deps: deps}
}

func (uc *deleteTherapistUseCase) Execute(ctx context.Context, therapistId string) error {
	if therapistId == "" {
		return fmt.Errorf("therapistId is empty")
	}

	therapist, err := uc.deps.TherapistRepo.GetById(ctx, therapistId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	if err := uc.deps.TherapistRepo.Delete(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDeletionFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
