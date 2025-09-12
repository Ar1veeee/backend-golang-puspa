package usecase

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	therapistErrors "backend-golang/internal/therapist/errors"
	"backend-golang/internal/therapist/repository"
	globalErrors "backend-golang/shared/errors"
	"context"
	"fmt"
)

type TherapistUseCase interface {
	CreateTherapistUseCase(ctx context.Context, req *dto.TherapistCreateRequest) error
	GetAllTherapistUseCase(ctx context.Context) ([]*dto.TherapistResponse, error)
}

type therapistUseCase struct {
	therapistRepo repository.TherapistRepository
	validator     TherapistValidator
	mapper        TherapistMapper
}

func NewTherapistUseCase(therapistRepo repository.TherapistRepository) TherapistUseCase {
	return &therapistUseCase{
		therapistRepo: therapistRepo,
		validator:     NewTherapistValidator(),
		mapper:        NewTherapistMapper(),
	}
}

func (uc *therapistUseCase) CreateTherapistUseCase(ctx context.Context, req *dto.TherapistCreateRequest) error {
	if err := uc.validator.validateCreateRequest(req); err != nil {
		return err
	}

	exists, err := uc.therapistRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrEmailExists
	}

	exists, err = uc.therapistRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrUsernameExists
	}

	tx := uc.therapistRepo.BeginTransaction(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", globalErrors.ErrDatabaseConnection)
	}

	user, therapist, err := uc.mapper.CreateRequestToUserAndTherapist(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := uc.therapistRepo.CreateUserWithTx(ctx, tx, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", therapistErrors.ErrTherapistCreationFailed, err)
	}

	if err := uc.therapistRepo.CreateTherapistWithTx(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", therapistErrors.ErrTherapistCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}

	return nil
}

func (uc *therapistUseCase) GetAllTherapistUseCase(ctx context.Context) ([]*dto.TherapistResponse, error) {
	therapists, err := uc.therapistRepo.GetAllTherapist(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", therapistErrors.ErrTherapistRetrievalFailed, err)
	}

	if therapists == nil {
		return []*dto.TherapistResponse{}, nil
	}

	responses := make([]*dto.TherapistResponse, 0, len(therapists))
	for _, therapist := range therapists {
		response := uc.mapper.AllTherapistsResponse(therapist.User, therapist)
		responses = append(responses, response)
	}

	return responses, nil
}
