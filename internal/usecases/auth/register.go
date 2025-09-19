package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type registerUseCase struct {
	deps *Dependencies
}

func NewRegisterUseCase(deps *Dependencies) RegisterUseCase {
	return &registerUseCase{deps: deps}
}

func (uc *registerUseCase) Execute(ctx context.Context, req *dto.RegisterRequest) error {
	if err := uc.deps.Validator.ValidateRegisterRequest(req); err != nil {
		log.Warn().Err(err).Str("email", req.Email).Msg("Registration validation failed")
		return err
	}

	tx := uc.deps.TxRepo.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	emailExists, usernameExists, err := uc.deps.UserRepo.CheckExisting(ctx, req.Email, req.Username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	if emailExists {
		tx.Rollback()
		return errors.ErrEmailExists
	}

	if usernameExists {
		tx.Rollback()
		return errors.ErrUsernameExists
	}

	parent, err := uc.deps.ParentRepo.GetByTempEmail(ctx, req.Email)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get parent of user")
		return errors.ErrInternalServer
	}
	if parent == nil {
		tx.Rollback()
		log.Warn().Str("email", req.Email).Msg("No parent found with temp_email")
		return errors.ErrEmailNotRegistered
	}

	user, err := uc.deps.Mapper.RegisterRequestToUser(req)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user entity")
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.deps.UserRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user in database")
		return errors.ErrInternalServer
	}

	if err := uc.deps.ParentRepo.UpdateUserId(ctx, tx, parent.TempEmail, user.Id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to update parent of user")
		return errors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: commit failed: %v", errors.ErrDatabaseConnection, err)
	}

	verificationToken, err := uc.deps.Mapper.CreateVerificationToken(user.Id)
	if err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("Failed to create verification token")
	}

	if verificationToken != nil {
		if err := uc.deps.VerifyTokenRepo.Create(ctx, verificationToken); err != nil {
			log.Error().Err(err).Str("userId", user.Id).Msg("Failed to save verification code")
			return errors.ErrInternalServer
		}

		verifyLink := fmt.Sprintf("http://localhost:3000/api/v1/auth/verify-account?token=%s", verificationToken.Token)
		if err := uc.deps.EmailService.SendVerificationEmail(user.Email, user.Username, verifyLink); err != nil {
			log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
			return errors.ErrInternalServer
		}
	}

	log.Info().Str("userId", user.Id).Str("email", req.Email).Msg("User registered and verification email sent")

	return nil
}
