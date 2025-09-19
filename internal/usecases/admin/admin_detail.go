package admin

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findAdminDetailUseCase struct {
	deps *Dependencies
}

func NewFindAdminDetailUseCase(deps *Dependencies) FindAdminDetailUseCase {
	return &findAdminDetailUseCase{deps: deps}
}

func (uc *findAdminDetailUseCase) Execute(ctx context.Context, adminId string) (*dto.AdminResponse, error) {
	if adminId == "" {
		return nil, fmt.Errorf("adminId is required")
	}

	adminDetail, err := uc.deps.AdminRepo.GetById(ctx, adminId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if adminDetail == nil {
		return nil, fmt.Errorf("adminDetail is nil")
	}

	var user *entities.User
	if adminDetail.User != nil {
		user = adminDetail.User
	}

	response, err := uc.deps.Mapper.AdminsResponse(user, adminDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to map admin %d: %w", adminDetail.Id, err)
	}

	return response, nil
}
