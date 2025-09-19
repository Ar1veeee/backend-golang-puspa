package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type FindPendingObservationsUseCase interface {
	Execute(ctx context.Context) ([]*dto.ObservationsResponse, error)
}

type FindScheduledObservationsUseCase interface {
	Execute(ctx context.Context) ([]*dto.ObservationsResponse, error)
}

type FindCompletedObservationsUseCase interface {
	Execute(ctx context.Context) ([]*dto.ObservationsResponse, error)
}

type FindObservationDetailUseCase interface {
	Execute(ctx context.Context, observationId int) (*dto.DetailObservationResponse, error)
}

type UpdateObservationDateUseCase interface {
	Execute(ctx context.Context, observationId int, req *dto.UpdateObservationDateRequest) error
}

type QuestionsUseCase interface {
	Execute(ctx context.Context, observationId int) ([]*dto.ObservationQuestionsResponse, error)
}

type SubmitObservationUseCase interface {
	Execute(ctx context.Context, observationId int, req *dto.SubmitObservationRequest) error
}
