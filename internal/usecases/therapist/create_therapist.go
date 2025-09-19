package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type createTherapistUseCase struct {
	deps *Dependencies
}

func NewCreateTherapistUseCase(deps *Dependencies) CreateTherapistUseCase {
	return &createTherapistUseCase{deps: deps}
}

func (uc *createTherapistUseCase) Execute(ctx context.Context, req *dto.TherapistCreateRequest) error {
	if err := uc.deps.Validator.ValidateCreateRequest(req); err != nil {
		return err
	}

	emailExists, usernameExists, err := uc.deps.UserRepo.CheckExisting(ctx, req.Email, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if emailExists {
		return errors.ErrEmailExists
	}

	if usernameExists {
		return errors.ErrUsernameExists
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, therapist, err := uc.deps.Mapper.CreateRequestToUserAndTherapist(req)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.deps.UserRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.deps.TherapistRepo.Create(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
