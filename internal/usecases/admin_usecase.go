package usecases

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/errors"
	"backend-golang/internal/mapper"
	"backend-golang/internal/validator"
	"context"
	"fmt"
)

type AdminUseCase interface {
	CreateAdminUseCase(ctx context.Context, req *dto.AdminCreateRequest) error
	FindAdminsUseCase(ctx context.Context) ([]*dto.AdminResponse, error)
	FindAdminDetailUseCase(ctx context.Context, adminId string) (*dto.AdminResponse, error)
	UpdateAdminUseCase(ctx context.Context, adminId string, req *dto.AdminUpdateRequest) error
	DeleteAdminWithTx(ctx context.Context, adminId string) error
}

type adminUseCase struct {
	txRepo    repositories.TransactionRepository
	userRepo  repositories.UserRepository
	adminRepo repositories.AdminRepository
	validator validator.AdminValidator
	mapper    mapper.AdminMapper
}

func NewAdminUseCase(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	adminRepo repositories.AdminRepository,
) AdminUseCase {
	return &adminUseCase{
		txRepo:    txRepo,
		userRepo:  userRepo,
		adminRepo: adminRepo,
		validator: validator.NewAdminValidator(),
		mapper:    mapper.NewAdminMapper(),
	}
}

func (uc *adminUseCase) CreateAdminUseCase(ctx context.Context, req *dto.AdminCreateRequest) error {
	if err := uc.validator.ValidateCreateRequest(req); err != nil {
		return err
	}

	emailExists, usernameExists, err := uc.userRepo.CheckExisting(ctx, req.Email, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if emailExists {
		return errors.ErrEmailExists
	}

	if usernameExists {
		return errors.ErrUsernameExists
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	user, admin, err := uc.mapper.CreateRequestToUserAndAdmin(req)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.userRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.adminRepo.Create(ctx, tx, admin); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}

func (uc *adminUseCase) FindAdminsUseCase(ctx context.Context) ([]*dto.AdminResponse, error) {
	admins, err := uc.adminRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if admins == nil {
		return []*dto.AdminResponse{}, nil
	}

	responses := make([]*dto.AdminResponse, 0, len(admins))
	for _, admin := range admins {
		if admin.User == nil {
			continue
		}

		response, err := uc.mapper.AdminsResponse(admin.User, admin)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %d: %w", admin.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (uc *adminUseCase) FindAdminDetailUseCase(ctx context.Context, adminId string) (*dto.AdminResponse, error) {
	if adminId == "" {
		return nil, fmt.Errorf("adminId is required")
	}

	adminDetail, err := uc.adminRepo.GetById(ctx, adminId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if adminDetail == nil {
		return nil, fmt.Errorf("adminDetail is nil")
	}

	var user *entities.User

	if adminDetail.User != nil {
		user = adminDetail.User
	}

	response, err := uc.mapper.AdminsResponse(
		user,
		adminDetail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map observation %d: %w", adminDetail.Id, err)
	}

	return response, nil
}

func (uc *adminUseCase) UpdateAdminUseCase(ctx context.Context, adminId string, req *dto.AdminUpdateRequest) error {
	if adminId == "" {
		return fmt.Errorf("admin id is required")
	}

	if err := uc.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	existingAdmin, err := uc.adminRepo.GetById(ctx, adminId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	if req.Email != "" && req.Email != existingAdmin.User.Email {
		emailExists, _, err := uc.userRepo.CheckExisting(ctx, req.Email, "")
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if emailExists {
			return errors.ErrEmailExists
		}
	}

	if req.Username != "" && req.Username != existingAdmin.User.Username {
		_, usernameExists, err := uc.userRepo.CheckExisting(ctx, "", req.Username)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if usernameExists {
			return errors.ErrUsernameExists
		}
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	updatedUser, updatedAdmin, err := uc.mapper.UpdateRequestToUserAndAdmin(req, existingAdmin)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.userRepo.Update(ctx, tx, updatedUser); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := uc.adminRepo.Update(ctx, tx, updatedAdmin); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}

func (uc *adminUseCase) DeleteAdminWithTx(ctx context.Context, adminId string) error {
	if adminId == "" {
		return fmt.Errorf("adminId is empty")
	}

	admin, err := uc.adminRepo.GetById(ctx, adminId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	if err := uc.adminRepo.Delete(ctx, tx, admin); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDeletionFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
