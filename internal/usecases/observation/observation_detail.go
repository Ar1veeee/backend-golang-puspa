package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findObservationDetailUseCase struct {
	deps *Dependencies
}

func NewFindObservationDetailUseCase(deps *Dependencies) FindObservationDetailUseCase {
	return &findObservationDetailUseCase{deps: deps}
}

func (uc *findObservationDetailUseCase) Execute(ctx context.Context, observationId int) (*dto.DetailObservationResponse, error) {
	if observationId == 0 {
		return nil, fmt.Errorf("observationId is required")
	}

	observationDetail, err := uc.deps.ObservationRepo.GetById(ctx, observationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if observationDetail == nil {
		return &dto.DetailObservationResponse{}, nil
	}

	var parent *entities.Parent
	var parentDetail *entities.ParentDetail
	var child *entities.Children

	if observationDetail.Children != nil {
		child = observationDetail.Children
		if child.Parent != nil {
			parent = child.Parent
			if len(parent.ParentDetail) > 0 {
				parentDetail = &parent.ParentDetail[0]
			}
		}
	}

	response, err := uc.deps.Mapper.ObservationDetailResponse(
		parent,
		parentDetail,
		child,
		observationDetail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map observation %d: %w", observationDetail.Id, err)
	}

	return response, nil
}
