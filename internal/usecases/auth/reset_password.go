package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"

	"github.com/rs/zerolog/log"
)

type resetPasswordUseCase struct {
	deps *Dependencies
}

func NewResetPasswordUseCase(deps *Dependencies) ResetPasswordUseCase {
	return &resetPasswordUseCase{deps: deps}
}

func (uc *resetPasswordUseCase) Execute(ctx context.Context, req *dto.ResetPasswordRequest) error {
	if err := uc.deps.Validator.ValidateResetPasswordRequest(req); err != nil {
		log.Warn().Err(err).Msg("Reset password validation failed")
		return err
	}

	token := req.Token

	if token == "" {
		log.Warn().Msg("Verification token is empty")
		return errors.ErrInvalidToken
	}

	verificationToken, err := uc.deps.VerifyTokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		log.Warn().Err(err).Str("code", req.Token).Msg("Invalid verification code")
		return errors.ErrInvalidToken
	}

	if verificationToken.IsExpired() {
		log.Warn().Str("code", req.Token).Str("userId", verificationToken.UserId).Msg("Verification code expired")
		return errors.ErrTokenExpired
	}

	user, err := uc.deps.Mapper.ResetPasswordRequestToUser(req)
	if err != nil {
		log.Warn().Err(err).Str("userId", user.Id).Msg("Failed to create forget password code")
		return errors.ErrInternalServer
	}

	if req.Password != req.ConfirmPassword {
		return errors.ErrPasswordNotSame
	}

	if err := uc.deps.UserRepo.UpdatePassword(ctx, verificationToken.UserId, user.Password); err != nil {
		log.Error().Err(err).Str("userId", verificationToken.UserId).Msg("Failed to update password")
		return errors.ErrInternalServer
	}

	log.Info().
		Str("userId", verificationToken.UserId).
		Msg("Update password successfully")

	return nil
}
