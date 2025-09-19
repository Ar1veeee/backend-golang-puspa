package admin

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findAdminsUseCase struct {
	deps *Dependencies
}

func NewFindAdminsUseCase(deps *Dependencies) FindAdminsUseCase {
	return &findAdminsUseCase{deps: deps}
}

func (uc *findAdminsUseCase) Execute(ctx context.Context) ([]*dto.AdminResponse, error) {
	admins, err := uc.deps.AdminRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if admins == nil {
		return []*dto.AdminResponse{}, nil
	}

	responses := make([]*dto.AdminResponse, 0, len(admins))
	for _, admin := range admins {
		if admin.User == nil {
			continue
		}

		response, err := uc.deps.Mapper.AdminsResponse(admin.User, admin)
		if err != nil {
			return nil, fmt.Errorf("failed to map admin %d: %w", admin.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
