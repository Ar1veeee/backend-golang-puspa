package child

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type FindChildUseCase interface {
	Execute(ctx context.Context) ([]*dto.ChildResponse, error)
}
