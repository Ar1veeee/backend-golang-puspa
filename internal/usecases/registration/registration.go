package registration

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type registrationUseCase struct {
	deps *Dependencies
}

func NewRegistrationUseCase(deps *Dependencies) RegistrationUseCase {
	return &registrationUseCase{deps: deps}
}

func (uc *registrationUseCase) Execute(ctx context.Context, req *dto.RegistrationRequest) error {
	if err := uc.deps.Validator.ValidateRegisterRequest(req); err != nil {
		return err
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	exists, err := uc.deps.ParentRepo.ExistByTempEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}
	if exists {
		tx.Rollback()
		return errors.ErrEmailExists
	}

	parent, parentDetail, child, observation, err := uc.deps.Mapper.CreateRequestToRegistration(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := uc.deps.ParentRepo.Create(ctx, tx, parent); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.deps.ParentDetailRepo.Create(ctx, tx, parentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err = uc.deps.ChildRepo.Create(ctx, tx, child); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.deps.ObservationRepo.Create(ctx, tx, observation); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
