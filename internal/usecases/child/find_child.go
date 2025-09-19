package child

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type findChildUseCase struct {
	deps *Dependencies
}

func NewFindChildUseCase(deps *Dependencies) FindChildUseCase {
	return &findChildUseCase{deps: deps}
}

func (uc *findChildUseCase) Execute(ctx context.Context) ([]*dto.ChildResponse, error) {
	childs, err := uc.deps.ChildRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	if childs == nil {
		return []*dto.ChildResponse{}, nil
	}

	responses := make([]*dto.ChildResponse, 0, len(childs))
	for _, child := range childs {
		var parentDetail *entities.ParentDetail
		if child.Parent != nil && len(child.Parent.ParentDetail) > 0 {
			parentDetail = &child.Parent.ParentDetail[0]
		}

		response, err := uc.deps.Mapper.ChildResponse(parentDetail, child)
		if err != nil {
			return nil, fmt.Errorf("failed to map child %d: %w", child.Id, err)
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
