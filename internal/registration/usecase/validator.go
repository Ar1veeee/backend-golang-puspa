package usecase

import (
	"backend-golang/internal/registration/delivery/http/dto"
	"backend-golang/shared/helpers"
)

type RegistrationValidator interface {
	ValidateRegisterRequest(req *dto.RegistrationRequest) error
}

type registrationValidator struct{}

func NewRegistrationValidator() RegistrationValidator {
	return &registrationValidator{}
}

func (v *registrationValidator) ValidateRegisterRequest(req *dto.RegistrationRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
