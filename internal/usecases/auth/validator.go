package auth

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"backend-golang/internal/validator"
)

type Validator interface {
	ValidateRegisterRequest(req *dto.RegisterRequest) error
	ValidateResendEmailRequest(req *dto.ResendTokenRequest) error
	ValidateResetPasswordRequest(req *dto.ResetPasswordRequest) error
	VerificationAccountRequest(req *dto.VerifyTokenRequest) error
	ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error
	ValidateLoginRequest(req *dto.LoginRequest) error
}

type authValidator struct{}

func NewAuthValidator() Validator {
	return &authValidator{}
}

func (v *authValidator) ValidateRegisterRequest(req *dto.RegisterRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return validator.IsValidPassword(req.Password)
}

func (v *authValidator) ValidateResendEmailRequest(req *dto.ResendTokenRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

func (v *authValidator) VerificationAccountRequest(req *dto.VerifyTokenRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	if req.Token == "" {
		return errors.ErrInvalidToken
	}

	return nil
}

func (v *authValidator) ValidateResetPasswordRequest(req *dto.ResetPasswordRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	if req.Token == "" {
		return errors.ErrInvalidToken
	}

	if req.Password != req.ConfirmPassword {
		return errors.ErrPasswordNotSame
	}

	return validator.IsValidPassword(req.Password)
}

func (v *authValidator) ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

func (v *authValidator) ValidateLoginRequest(req *dto.LoginRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
