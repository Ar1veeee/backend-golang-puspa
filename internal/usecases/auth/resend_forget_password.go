package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/errors"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type resendForgetPasswordUseCase struct {
	deps *Dependencies
}

func NewResendForgetPasswordUseCase(deps *Dependencies) ResendForgetPasswordUseCase {
	return &resendForgetPasswordUseCase{deps: deps}
}

func (uc *resendForgetPasswordUseCase) Execute(ctx context.Context, req *dto.ResendTokenRequest) error {
	if err := uc.deps.Validator.ValidateResendEmailRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	user, err := uc.deps.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to find user by email")
		return errors.ErrEmailNotFound
	}

	if user.IsActive {
		log.Warn().Str("email", req.Email).Msg("User is already active")
		return errors.ErrEmailAlreadyVerified
	}

	existingToken, err := uc.deps.VerifyTokenRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to check existing token")
		return errors.ErrInternalServer
	}

	var verificationToken *entities.VerificationToken
	var verifyLink string

	if existingToken != nil && !existingToken.IsExpired() {
		verificationToken = existingToken
		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", existingToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Reusing existing valid token")
	} else {
		verificationToken, err = uc.deps.Mapper.CreateVerificationToken(user.Id)
		if err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create new verification token")
			return errors.ErrInternalServer
		}

		if err := uc.deps.VerifyTokenRepo.Create(ctx, verificationToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save new verification token")
			return errors.ErrInternalServer
		}

		verifyLink = fmt.Sprintf("http://localhost:3000/api/v1/auth/update-password?token=%s", verificationToken.Token)
		log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Created and saved new verification token")
	}

	if err := uc.deps.EmailService.SendResetPasswordEmail(user.Email, user.Username, verifyLink); err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
		return errors.ErrInternalServer
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("Verification email resent successfully")
	return nil
}
