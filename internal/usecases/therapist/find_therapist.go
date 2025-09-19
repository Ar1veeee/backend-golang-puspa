package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findTherapistsUseCase struct {
	deps *Dependencies
}

func NewFindTherapistsUseCase(deps *Dependencies) FindTherapistsUseCase {
	return &findTherapistsUseCase{deps: deps}
}

func (uc *findTherapistsUseCase) Execute(ctx context.Context) ([]*dto.TherapistResponse, error) {
	therapists, err := uc.deps.TherapistRepo.GetAll(ctx)
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

		response, err := uc.deps.Mapper.TherapistsResponse(therapist.User, therapist)
		if err != nil {
			return nil, fmt.Errorf("failed to map observation %s: %w", therapist.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
