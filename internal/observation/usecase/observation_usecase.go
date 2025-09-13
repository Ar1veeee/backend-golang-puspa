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
	GetPendingObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error)
	GetCompletedObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error)
	GetObservationDetail(ctx context.Context, id int) (*dto.DetailObservationResponse, error)
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

func (uc *observationUseCase) GetPendingObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error) {
	observations, err := uc.observationRepo.GetPendingObservations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", observationErrors.ErrObservationRetrievalFailed, err)
	}

	if observations == nil {
		return []*dto.ObservationsResponse{}, nil
	}

	responses := make([]*dto.ObservationsResponse, 0, len(observations))
	for _, observation := range observations {
		if observation.Children == nil || observation.Children.Parent == nil {
			continue
		}

		var parentDetail *entity.ParentDetail
		if len(observation.Children.Parent.ParentDetail) > 0 {
			parentDetail = &observation.Children.Parent.ParentDetail[0]
		}

		response, err := uc.mapper.ObservationsResponse(parentDetail, observation.Children, observation)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %d: %w", observation.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (uc *observationUseCase) GetCompletedObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error) {
	observations, err := uc.observationRepo.GetCompletedObservations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", observationErrors.ErrObservationRetrievalFailed, err)
	}

	if observations == nil {
		return []*dto.ObservationsResponse{}, nil
	}

	responses := make([]*dto.ObservationsResponse, 0, len(observations))
	for _, observation := range observations {
		if observation.Children == nil || observation.Children.Parent == nil {
			continue
		}

		var parentDetail *entity.ParentDetail
		if len(observation.Children.Parent.ParentDetail) > 0 {
			parentDetail = &observation.Children.Parent.ParentDetail[0]
		}

		response, err := uc.mapper.ObservationsResponse(parentDetail, observation.Children, observation)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %d: %w", observation.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (uc *observationUseCase) GetObservationDetail(ctx context.Context, id int) (*dto.DetailObservationResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("observation id is required")
	}
	observationDetail, err := uc.observationRepo.GetObservationById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", observationErrors.ErrObservationRetrievalFailed, err)
	}

	if observationDetail == nil {
		return &dto.DetailObservationResponse{}, nil
	}

	var parent *entity.Parent
	var parentDetail *entity.ParentDetail
	var child *entity.Children

	if observationDetail.Children != nil {
		child = observationDetail.Children
		if child.Parent != nil {
			parent = child.Parent
			if len(parent.ParentDetail) > 0 {
				parentDetail = &parent.ParentDetail[0]
			}
		}
	}

	response, err := uc.mapper.ObservationDetailResponse(
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
