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

type TherapistUseCase interface {
	CreateTherapistUseCase(ctx context.Context, req *dto.TherapistCreateRequest) error
	FindTherapistsUseCase(ctx context.Context) ([]*dto.TherapistResponse, error)
	FindTherapistDetailUseCase(ctx context.Context, therapistId string) (*dto.TherapistResponse, error)
	UpdateTherapistUseCase(ctx context.Context, therapistId string, req *dto.TherapistUpdateRequest) error
	DeleteTherapistWithTx(ctx context.Context, therapistId string) error
}

type therapistUseCase struct {
	txRepo        repositories.TransactionRepository
	userRepo      repositories.UserRepository
	therapistRepo repositories.TherapistRepository
	validator     validator.TherapistValidator
	mapper        mapper.TherapistMapper
}

func NewTherapistUseCase(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	therapistRepo repositories.TherapistRepository,
) TherapistUseCase {
	return &therapistUseCase{
		txRepo:        txRepo,
		userRepo:      userRepo,
		therapistRepo: therapistRepo,
		validator:     validator.NewTherapistValidator(),
		mapper:        mapper.NewTherapistMapper(),
	}
}

func (uc *therapistUseCase) CreateTherapistUseCase(ctx context.Context, req *dto.TherapistCreateRequest) error {
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

	user, therapist, err := uc.mapper.CreateRequestToUserAndTherapist(req)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.userRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := uc.therapistRepo.Create(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}

func (uc *therapistUseCase) FindTherapistsUseCase(ctx context.Context) ([]*dto.TherapistResponse, error) {
	therapists, err := uc.therapistRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if therapists == nil {
		return []*dto.TherapistResponse{}, nil
	}

	responses := make([]*dto.TherapistResponse, 0, len(therapists))
	for _, therapist := range therapists {
		if therapist.User == nil {
			continue
		}

		response, err := uc.mapper.TherapistsResponse(therapist.User, therapist)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %s: %w", therapist.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (uc *therapistUseCase) FindTherapistDetailUseCase(ctx context.Context, therapistId string) (*dto.TherapistResponse, error) {
	if therapistId == "" {
		return nil, fmt.Errorf("therapistId is required")
	}

	therapistDetail, err := uc.therapistRepo.GetById(ctx, therapistId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if therapistDetail == nil {
		return nil, fmt.Errorf("therapistDetail is nil")
	}

	var user *entities.User

	if therapistDetail.User != nil {
		user = therapistDetail.User
	}

	response, err := uc.mapper.TherapistsResponse(
		user,
		therapistDetail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map observation %s: %w", therapistDetail.Id, err)
	}

	return response, nil
}

func (uc *therapistUseCase) UpdateTherapistUseCase(ctx context.Context, therapistId string, req *dto.TherapistUpdateRequest) error {
	if therapistId == "" {
		return fmt.Errorf("therapist id is required")
	}

	if err := uc.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	existingTherapist, err := uc.therapistRepo.GetById(ctx, therapistId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	if req.Email != "" && req.Email != existingTherapist.User.Email {
		emailExists, _, err := uc.userRepo.CheckExisting(ctx, req.Email, "")
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
		}
		if emailExists {
			return errors.ErrEmailExists
		}
	}

	if req.Username != "" && req.Username != existingTherapist.User.Username {
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

	updatedUser, updatedTherapist, err := uc.mapper.UpdateRequestToUserAndTherapist(req, existingTherapist)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if err := uc.userRepo.Update(ctx, tx, updatedUser); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := uc.therapistRepo.Update(ctx, tx, updatedTherapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrUpdateFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}

func (uc *therapistUseCase) DeleteTherapistWithTx(ctx context.Context, therapistId string) error {
	if therapistId == "" {
		return fmt.Errorf("therapistId is empty")
	}

	therapist, err := uc.therapistRepo.GetById(ctx, therapistId)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNotFound, err)
	}

	tx := uc.txRepo.Begin(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", errors.ErrDatabaseConnection)
	}

	if err := uc.therapistRepo.Delete(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDeletionFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	return nil
}
