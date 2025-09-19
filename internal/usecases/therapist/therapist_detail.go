package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findTherapistDetailUseCase struct {
	deps *Dependencies
}

func NewFindTherapistDetailUseCase(deps *Dependencies) FindTherapistDetailUseCase {
	return &findTherapistDetailUseCase{deps: deps}
}

func (uc *findTherapistDetailUseCase) Execute(ctx context.Context, therapistId string) (*dto.TherapistResponse, error) {
	if therapistId == "" {
		return nil, fmt.Errorf("therapistId is required")
	}

	therapistDetail, err := uc.deps.TherapistRepo.GetById(ctx, therapistId)
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

	response, err := uc.deps.Mapper.TherapistsResponse(
		user,
		therapistDetail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map observation %s: %w", therapistDetail.Id, err)
	}

	return response, nil
}
