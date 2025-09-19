package therapist

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/validator"
)

type Validator interface {
	ValidateCreateRequest(req *dto.TherapistCreateRequest) error
	ValidateUpdateRequest(req *dto.TherapistUpdateRequest) error
}

type therapistValidator struct{}

func NewTherapistValidator() Validator {
	return &therapistValidator{}
}

func (v *therapistValidator) ValidateCreateRequest(req *dto.TherapistCreateRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return validator.IsValidPassword(req.Password)
}

func (v *therapistValidator) ValidateUpdateRequest(req *dto.TherapistUpdateRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return validator.IsValidPassword(req.Password)
}
