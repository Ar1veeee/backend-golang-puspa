package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findPendingObservationsUseCase struct {
	deps *Dependencies
}

func NewFindPendingObservationsUseCase(deps *Dependencies) FindPendingObservationsUseCase {
	return &findPendingObservationsUseCase{deps: deps}
}

func (uc *findPendingObservationsUseCase) Execute(ctx context.Context) ([]*dto.ObservationsResponse, error) {
	observations, err := uc.deps.ObservationRepo.GetByPendingStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if observations == nil {
		return []*dto.ObservationsResponse{}, nil
	}

	responses := make([]*dto.ObservationsResponse, 0, len(observations))
	for _, observation := range observations {
		if observation.Children == nil || observation.Children.Parent == nil {
			continue
		}

		var parentDetail *entities.ParentDetail
		if len(observation.Children.Parent.ParentDetail) > 0 {
			parentDetail = &observation.Children.Parent.ParentDetail[0]
		}

		response, err := uc.deps.Mapper.ObservationsResponse(parentDetail, observation.Children, observation)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %d: %w", observation.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
