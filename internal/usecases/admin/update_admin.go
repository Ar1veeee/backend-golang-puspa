package admin

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type updateAdminUseCase struct {
	deps *Dependencies
}

func NewUpdateAdminUseCase(deps *Dependencies) UpdateAdminUseCase {
	return &updateAdminUseCase{deps: deps}
}

func (uc *updateAdminUseCase) Execute(ctx context.Context, adminId string, req *dto.AdminUpdateRequest) error {
	if adminId == "" {
		return fmt.Errorf("admin id is required")
	}

	if err := uc.deps.Validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	existingAdmin, err := uc.deps.AdminRepo.GetById(ctx, adminId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	if req.Email != "" && req.Email != existingAdmin.User.Email {
		emailExists, _, err := uc.deps.UserRepo.CheckExisting(ctx, req.Email, "")
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if emailExists {
			return errors.ErrEmailExists
		}
	}

	if req.Username != "" && req.Username != existingAdmin.User.Username {
		_, usernameExists, err := uc.deps.UserRepo.CheckExisting(ctx, "", req.Username)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if usernameExists {
			return errors.ErrUsernameExists
		}
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	updatedUser, updatedAdmin, err := uc.deps.Mapper.UpdateRequestToUserAndAdmin(req, existingAdmin)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.deps.UserRepo.Update(ctx, tx, updatedUser); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := uc.deps.AdminRepo.Update(ctx, tx, updatedAdmin); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
