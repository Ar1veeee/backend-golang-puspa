package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type updateTherapistUseCase struct {
	deps *Dependencies
}

func NewUpdateTherapistUseCase(deps *Dependencies) UpdateTherapistUseCase {
	return &updateTherapistUseCase{deps: deps}
}

func (uc *updateTherapistUseCase) Execute(ctx context.Context, therapistId string, req *dto.TherapistUpdateRequest) error {
	if err := uc.deps.Validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	existingTherapist, err := uc.deps.TherapistRepo.GetById(ctx, therapistId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	if req.Email != "" && req.Email != existingTherapist.User.Email {
		emailExists, _, err := uc.deps.UserRepo.CheckExisting(ctx, req.Email, "")
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if emailExists {
			return errors.ErrEmailExists
		}
	}

	if req.Username != "" && req.Username != existingTherapist.User.Username {
		_, usernameExists, err := uc.deps.UserRepo.CheckExisting(ctx, "", req.Username)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if usernameExists {
			return errors.ErrUsernameExists
		}
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updatedUser, updatedTherapist, err := uc.deps.Mapper.UpdateRequestToUserAndTherapist(req, existingTherapist)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.deps.UserRepo.Update(ctx, tx, updatedUser); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := uc.deps.TherapistRepo.Update(ctx, tx, updatedTherapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
