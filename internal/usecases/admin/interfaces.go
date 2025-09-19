package admin

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type CreateAdminUseCase interface {
	Execute(ctx context.Context, req *dto.AdminCreateRequest) error
}

type FindAdminsUseCase interface {
	Execute(ctx context.Context) ([]*dto.AdminResponse, error)
}

type FindAdminDetailUseCase interface {
	Execute(ctx context.Context, adminId string) (*dto.AdminResponse, error)
}

type UpdateAdminUseCase interface {
	Execute(ctx context.Context, adminId string, req *dto.AdminUpdateRequest) error
}

type DeleteAdminUseCase interface {
	Execute(ctx context.Context, adminId string) error
}
