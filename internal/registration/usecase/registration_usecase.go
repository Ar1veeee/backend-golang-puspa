package usecase

import (
	"backend-golang/internal/registration/delivery/http/dto"
	registrationErrors "backend-golang/internal/registration/errors"
	"backend-golang/internal/registration/repository"
	globalErrors "backend-golang/shared/errors"
	"context"
	"fmt"
)

type RegistrationUseCase interface {
	RegistrationUseCase(ctx context.Context, req *dto.RegistrationRequest) error
}

type registrationUseCase struct {
	registrationRepo repository.RegistrationRepository
	validator        RegistrationValidator
	mapper           RegistrationMapper
}

func NewRegistrationUseCase(
	registrationRepo repository.RegistrationRepository,
) RegistrationUseCase {
	return &registrationUseCase{
		registrationRepo: registrationRepo,
		validator:        NewRegistrationValidator(),
		mapper:           NewRegistrationMapper(),
	}
}

func (uc *registrationUseCase) RegistrationUseCase(ctx context.Context, req *dto.RegistrationRequest) error {
	if err := uc.validator.ValidateRegisterRequest(req); err != nil {
		return err
	}

	tx := uc.registrationRepo.BeginTransaction(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", globalErrors.ErrDatabaseConnection)
	}

	exists, err := uc.registrationRepo.ExistsByEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		tx.Rollback()
		return globalErrors.ErrEmailExists
	}

	parent, parentDetail, child, observation, err := uc.mapper.CreateRequestToRegistration(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := uc.registrationRepo.CreateParentWithTx(ctx, tx, parent); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := uc.registrationRepo.CreateParentDetailWithTx(ctx, tx, parentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err = uc.registrationRepo.CreateChildWithTx(ctx, tx, child); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := uc.registrationRepo.CreateObservationWithTx(ctx, tx, observation); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}

	return nil
}
