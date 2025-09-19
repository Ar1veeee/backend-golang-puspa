package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type forgetPasswordUseCase struct {
	deps *Dependencies
}

func NewForgetPasswordUseCase(deps *Dependencies) ForgetPasswordUseCase {
	return &forgetPasswordUseCase{deps: deps}
}

func (uc *forgetPasswordUseCase) Execute(ctx context.Context, req *dto.ForgetPasswordRequest) error {
	if err := uc.deps.Validator.ValidateForgetPasswordRequest(req); err != nil {
		return err
	}

	user, err := uc.deps.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Reset password attempt failed: email not found")
		return errors.ErrEmailNotFound
	}

	verificationCode, err := uc.deps.Mapper.CreateVerificationToken(user.Id)
	if err != nil {
		log.Warn().Err(err).Str("userId", user.Id).Msg("Failed to create forget password code")
		return errors.ErrInternalServer
	}

	if verificationCode != nil {
		if err := uc.deps.VerifyTokenRepo.Create(ctx, verificationCode); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
			return errors.ErrInternalServer
		}

		verifyLink := fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", verificationCode.Token)
		if err := uc.deps.EmailService.SendResetPasswordEmail(user.Email, user.Username, verifyLink); err != nil {
			log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
			return errors.ErrInternalServer
		}
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Forget password code sent")
	return nil
}
