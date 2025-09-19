package admin

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/validator"
	"fmt"
)

type Validator interface {
	ValidateCreateRequest(req *dto.AdminCreateRequest) error
	ValidateUpdateRequest(req *dto.AdminUpdateRequest) error
}

type adminValidator struct{}

func NewAdminValidator() Validator {
	return &adminValidator{}
}

func (a *adminValidator) ValidateCreateRequest(req *dto.AdminCreateRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return validator.IsValidPassword(req.Password)
}

func (v *adminValidator) ValidateUpdateRequest(req *dto.AdminUpdateRequest) error {
	if req.Username == "" && req.Email == "" && req.Password == "" &&
		req.AdminName == "" && req.AdminPhone == "" {
		return fmt.Errorf("at least one field must be provided for update")
	}

	if req.Password != "" {
		return validator.IsValidPassword(req.Password)
	}

	return nil
}
