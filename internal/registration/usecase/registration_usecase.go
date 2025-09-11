package usecase

import (
	"backend-golang/internal/registration/delivery/http/dto"
	registrationErrors "backend-golang/internal/registration/errors"
	"backend-golang/internal/registration/repository"
	globalErrors "backend-golang/shared/errors"
	"context"
	"fmt"
)

type RegistrationService interface {
	Registration(ctx context.Context, req *dto.RegistrationRequest) error
}

type registrationService struct {
	registrationRepo repository.RegistrationRepository
	validator        RegistrationValidator
	mapper           RegistrationMapper
}

func NewRegistrationService(
	registrationRepo repository.RegistrationRepository,
) RegistrationService {
	return &registrationService{
		registrationRepo: registrationRepo,
		validator:        NewRegistrationValidator(),
		mapper:           NewRegistrationMapper(),
	}
}

func (s *registrationService) Registration(ctx context.Context, req *dto.RegistrationRequest) error {
	if err := s.validator.validateRegisterRequest(req); err != nil {
		return err
	}

	tx := s.registrationRepo.BeginTransaction(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", globalErrors.ErrDatabaseConnection)
	}

	exists, err := s.registrationRepo.ExistsByEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		tx.Rollback()
		return globalErrors.ErrEmailExists
	}

	parent, parentDetail, child, observation, err := s.mapper.createRequestToRegistration(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.registrationRepo.CreateParentWithTx(ctx, tx, parent); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateParentDetailWithTx(ctx, tx, parentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateChildWithTx(ctx, tx, child); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateObservationWithTx(ctx, tx, observation); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}

	return nil
}
