package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type submitObservationUseCase struct {
	deps *Dependencies
}

func NewSubmitObservationUseCase(deps *Dependencies) SubmitObservationUseCase {
	return &submitObservationUseCase{deps: deps}
}

func (uc *submitObservationUseCase) Execute(ctx context.Context, observationId int, req *dto.SubmitObservationRequest) error {
	if observationId == 0 {
		return fmt.Errorf("ObservationId is required")
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	observation, answers, err := uc.deps.Mapper.UpdateToObservationAndCreateToAnswer(ctx, observationId, req)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.deps.ObservationRepo.UpdateAfterObservation(
		ctx,
		tx,
		observationId,
		observation.TherapistId,
		observation.TotalScore,
		observation.Conclusion,
		observation.Recommendation,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := uc.deps.ObservationAnswerRepo.Create(ctx, tx, answers); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
