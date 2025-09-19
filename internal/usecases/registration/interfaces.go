package registration

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type RegistrationUseCase interface {
	Execute(ctx context.Context, req *dto.RegistrationRequest) error
}
