package usecase

import (
	"backend-golang/internal/observation/delivery/http/dto"
	"backend-golang/internal/observation/entity"
	observationErrors "backend-golang/internal/observation/errors"
	"backend-golang/internal/observation/repository"
	"context"
	"fmt"
)

type ObservationUseCase interface {
	GetAllObservationsUseCase(ctx context.Context) ([]*dto.PendingObservationsResponse, error)
}

type observationUseCase struct {
	observationRepo repository.ObservationRepository
	mapper          ObservationMapper
}

func NewObservationUseCase(observationRepo repository.ObservationRepository) ObservationUseCase {
	return &observationUseCase{
		observationRepo: observationRepo,
		mapper:          NewObservationMapper(),
	}
}

func (uc *observationUseCase) GetAllObservationsUseCase(ctx context.Context) ([]*dto.PendingObservationsResponse, error) {
	observations, err := uc.observationRepo.GetPendingObservations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", observationErrors.ErrObservationRetrievalFailed, err)
	}

	if observations == nil {
		return []*dto.PendingObservationsResponse{}, nil
	}

	responses := make([]*dto.PendingObservationsResponse, 0, len(observations))
	for _, observation := range observations {
		if observation.Children == nil || observation.Children.Parent == nil {
			continue
		}

		var parentDetail *entity.ParentDetail
		if len(observation.Children.Parent.ParentDetail) > 0 {
			parentDetail = &observation.Children.Parent.ParentDetail[0]
		}

		response, err := uc.mapper.AllObservationsResponse(parentDetail, observation.Children, observation)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %d: %w", observation.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
