package usecase

import (
	"backend-golang/internal/auth/delivery/http/dto"
	"backend-golang/shared/helpers"
	"errors"
	"fmt"
	"strings"
)

type AuthValidator interface {
	ValidateRegisterRequest(req *dto.RegisterRequest) error
	ValidateVerifyCodeRequest(req *dto.VerifyCodeRequest) error
	ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error
	ValidateLoginRequest(req *dto.LoginRequest) error
}

type authValidator struct{}

func NewAuthValidator() AuthValidator {
	return &authValidator{}
}

func (v *authValidator) ValidateRegisterRequest(req *dto.RegisterRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return helpers.IsValidPassword(req.Password)
}

func (v *authValidator) ValidateVerifyCodeRequest(req *dto.VerifyCodeRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.Code == "" {
		return errors.New("verification code cannot be empty")
	}

	return nil
}

func (v *authValidator) ValidateForgetPasswordRequest(req *dto.ForgetPasswordRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if strings.Contains(req.Email, "@") {
		if !helpers.IsValidEmail(req.Email) {
			return fmt.Errorf("invalid email format")
		}
	}

	return nil
}

func (v *authValidator) ValidateLoginRequest(req *dto.LoginRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if strings.Contains(req.Identifier, "@") {
		if !helpers.IsValidEmail(req.Identifier) {
			return fmt.Errorf("invalid email format")
		}
	}

	return nil
}
