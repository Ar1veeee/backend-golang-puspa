package validator

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
)

type AuthValidator interface {
	ValidateRegisterRequest(req *dto.RegisterRequest) error
	ValidateResendEmailRequest(req *dto.ResendTokenRequest) error
	ValidateResetPasswordRequest(req *dto.ResetPasswordRequest) error
	VerificationAccountRequest(req *dto.VerifyTokenRequest) error
	ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error
	ValidateLoginRequest(req *dto.LoginRequest) error
}

type authValidator struct{}

func NewAuthValidator() AuthValidator {
	return &authValidator{}
}

func (v *authValidator) ValidateRegisterRequest(req *dto.RegisterRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return IsValidPassword(req.Password)
}

func (v *authValidator) ValidateResendEmailRequest(req *dto.ResendTokenRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

func (v *authValidator) VerificationAccountRequest(req *dto.VerifyTokenRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	if req.Token == "" {
		return errors.ErrInvalidToken
	}

	return nil
}

func (v *authValidator) ValidateResetPasswordRequest(req *dto.ResetPasswordRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	if req.Token == "" {
		return errors.ErrInvalidToken
	}

	if req.Password != req.ConfirmPassword {
		return errors.ErrPasswordNotSame
	}

	return IsValidPassword(req.Password)
}

func (v *authValidator) ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

func (v *authValidator) ValidateLoginRequest(req *dto.LoginRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
