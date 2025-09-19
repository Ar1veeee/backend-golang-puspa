package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type updateObservationDateUseCase struct {
	deps *Dependencies
}

func NewUpdateObservationDateUseCase(deps *Dependencies) UpdateObservationDateUseCase {
	return &updateObservationDateUseCase{deps: deps}
}

func (uc *updateObservationDateUseCase) Execute(ctx context.Context, observationId int, req *dto.UpdateObservationDateRequest) error {
	if err := uc.deps.Validator.ValidateUpdateScheduledDateRequest(req); err != nil {
		return err
	}

	if err := uc.deps.ObservationRepo.UpdateScheduledDate(ctx, observationId, req.ScheduledDate); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	return nil
}
