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

type ObservationUseCase interface {
	FindPendingObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error)
	FindCompletedObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error)
	FindObservationDetailUseCase(ctx context.Context, id int) (*dto.DetailObservationResponse, error)
}

type observationUseCase struct {
	observationRepo repositories.ObservationRepository
	validator       validator.ObservationValidator
	mapper          mapper.ObservationMapper
}

func NewObservationUseCase(
	observationRepo repositories.ObservationRepository,
) ObservationUseCase {
	return &observationUseCase{
		observationRepo: observationRepo,
		mapper:          mapper.NewObservationMapper(),
	}
}

func (uc *observationUseCase) FindPendingObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error) {
	observations, err := uc.observationRepo.GetByPendingStatus(ctx)
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

func (uc *observationUseCase) FindCompletedObservationsUseCase(ctx context.Context) ([]*dto.ObservationsResponse, error) {
	observations, err := uc.observationRepo.GetByCompletedStatus(ctx)
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

func (uc *observationUseCase) FindObservationDetailUseCase(ctx context.Context, observationId int) (*dto.DetailObservationResponse, error) {
	if observationId == 0 {
		return nil, fmt.Errorf("observationId is required")
	}

	observationDetail, err := uc.observationRepo.GetById(ctx, observationId)
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
