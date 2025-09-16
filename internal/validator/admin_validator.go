package validator

import (
	"backend-golang/internal/adapters/http/dto"
	"fmt"
)

type AdminValidator interface {
	ValidateCreateRequest(req *dto.AdminCreateRequest) error
	ValidateUpdateRequest(req *dto.AdminUpdateRequest) error
}

type adminValidator struct{}

func NewAdminValidator() AdminValidator {
	return &adminValidator{}
}

func (a *adminValidator) ValidateCreateRequest(req *dto.AdminCreateRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return IsValidPassword(req.Password)
}

func (v *adminValidator) ValidateUpdateRequest(req *dto.AdminUpdateRequest) error {
	if req.Username == "" && req.Email == "" && req.Password == "" &&
		req.AdminName == "" && req.AdminPhone == "" {
		return fmt.Errorf("at least one field must be provided for update")
	}

	if req.Password != "" {
		return IsValidPassword(req.Password)
	}

	return nil
}
