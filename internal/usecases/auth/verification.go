package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"

	"github.com/rs/zerolog/log"
)

type verificationAccountUseCase struct {
	deps *Dependencies
}

func NewVerificationAccountUseCase(deps *Dependencies) VerificationAccountUseCase {
	return &verificationAccountUseCase{deps: deps}
}

func (uc *verificationAccountUseCase) Execute(ctx context.Context, req *dto.VerifyTokenRequest) error {
	if err := uc.deps.Validator.VerificationAccountRequest(req); err != nil {
		log.Warn().Err(err).Str("token", req.Token).Msg("Verification account failed")
		return err
	}

	verificationToken, err := uc.deps.VerifyTokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		log.Warn().Err(err).Str("code", req.Token).Msg("Invalid verification code")
		return errors.ErrInvalidToken
	}

	if err := uc.deps.VerifyTokenRepo.UpdateStatus(ctx, verificationToken.Token); err != nil {
		log.Error().Err(err).Str("code", req.Token).Msg("Failed to update verification code")
		return errors.ErrInternalServer
	}

	if err := uc.deps.ParentRepo.UpdateRegistrationStatus(ctx, verificationToken.UserId); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to complete registration")
		return errors.ErrInternalServer
	}

	if err := uc.deps.UserRepo.UpdateActiveStatus(ctx, verificationToken.UserId, true); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to activate user")
		return errors.ErrInternalServer
	}

	log.Info().
		Str("userId", verificationToken.UserId).
		Msg("Email verified successfully")

	return nil
}
