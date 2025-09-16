package validator

import (
	"backend-golang/internal/adapters/http/dto"
	"fmt"
)

type TherapistValidator interface {
	ValidateCreateRequest(req *dto.TherapistCreateRequest) error
	ValidateUpdateRequest(req *dto.TherapistUpdateRequest) error
}

type therapistValidator struct{}

func NewTherapistValidator() TherapistValidator {
	return &therapistValidator{}
}

func (v *therapistValidator) ValidateCreateRequest(req *dto.TherapistCreateRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return IsValidPassword(req.Password)
}

func (v *therapistValidator) ValidateUpdateRequest(req *dto.TherapistUpdateRequest) error {
	if req.Username == "" && req.Email == "" && req.Password == "" &&
		req.TherapistName == "" && req.TherapistPhone == "" && req.TherapistSection == "" {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return IsValidPassword(req.Password)
}
