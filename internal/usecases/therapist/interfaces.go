package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type CreateTherapistUseCase interface {
	Execute(ctx context.Context, req *dto.TherapistCreateRequest) error
}

type FindTherapistsUseCase interface {
	Execute(ctx context.Context) ([]*dto.TherapistResponse, error)
}

type FindTherapistDetailUseCase interface {
	Execute(ctx context.Context, therapistId string) (*dto.TherapistResponse, error)
}

type UpdateTherapistUseCase interface {
	Execute(ctx context.Context, therapistId string, req *dto.TherapistUpdateRequest) error
}

type DeleteTherapistUseCase interface {
	Execute(ctx context.Context, therapistId string) error
}
