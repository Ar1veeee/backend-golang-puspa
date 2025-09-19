package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"context"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, req *dto.RegisterRequest) error
}

type LoginUseCase interface {
	Execute(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, req *dto.ResetPasswordRequest) error
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, refreshToken string) error
}

type ResendVerificationAccountUseCase interface {
	Execute(ctx context.Context, req *dto.ResendTokenRequest) error
}

type VerificationAccountUseCase interface {
	Execute(ctx context.Context, req *dto.VerifyTokenRequest) error
}

type ForgetPasswordUseCase interface {
	Execute(ctx context.Context, req *dto.ForgetPasswordRequest) error
}

type ResendForgetPasswordUseCase interface {
	Execute(ctx context.Context, req *dto.ResendTokenRequest) error
}
