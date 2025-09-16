package usecases

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/errors"
	"backend-golang/internal/mapper"
	"backend-golang/internal/validator"
	"context"
	"fmt"
)

type RegistrationUseCase interface {
	RegistrationUseCase(ctx context.Context, req *dto.RegistrationRequest) error
}

type registrationUseCase struct {
	txRepo           repositories.TransactionRepository
	parentRepo       repositories.ParentRepository
	parentDetailRepo repositories.ParentDetailRepository
	childRepo        repositories.ChildRepository
	observationRepo  repositories.ObservationRepository
	validator        validator.RegistrationValidator
	mapper           mapper.RegistrationMapper
}

func NewRegistrationUseCase(
	txRepo repositories.TransactionRepository,
	parentRepo repositories.ParentRepository,
	parentDetailRepo repositories.ParentDetailRepository,
	childRepo repositories.ChildRepository,
	observationRepo repositories.ObservationRepository,
) RegistrationUseCase {
	return &registrationUseCase{
		txRepo:           txRepo,
		parentRepo:       parentRepo,
		parentDetailRepo: parentDetailRepo,
		childRepo:        childRepo,
		observationRepo:  observationRepo,
		validator:        validator.NewRegistrationValidator(),
		mapper:           mapper.NewRegistrationMapper(),
	}
}

func (uc *registrationUseCase) RegistrationUseCase(ctx context.Context, req *dto.RegistrationRequest) error {
	if err := uc.validator.ValidateRegisterRequest(req); err != nil {
		return err
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	exists, err := uc.parentRepo.ExistByTempEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}
	if exists {
		tx.Rollback()
		return errors.ErrEmailExists
	}

	parent, parentDetail, child, observation, err := uc.mapper.CreateRequestToRegistration(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := uc.parentRepo.Create(ctx, tx, parent); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.parentDetailRepo.Create(ctx, tx, parentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err = uc.childRepo.Create(ctx, tx, child); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.observationRepo.Create(ctx, tx, observation); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
