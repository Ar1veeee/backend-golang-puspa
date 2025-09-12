package usecase

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	"backend-golang/shared/helpers"
)

type TherapistValidator interface {
	validateCreateRequest(req *dto.TherapistCreateRequest) error
	validateUpdateRequest(req *dto.TherapistUpdateRequest) error
}

type therapistValidator struct{}

func NewTherapistValidator() TherapistValidator {
	return &therapistValidator{}
}

func (v *therapistValidator) validateCreateRequest(req *dto.TherapistCreateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return helpers.IsValidPassword(req.Password)
}

func (v *therapistValidator) validateUpdateRequest(req *dto.TherapistUpdateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.Password != "" {
		return helpers.IsValidPassword(req.Password)
	}

	if req.Role != "" {

	}

	return nil
}
