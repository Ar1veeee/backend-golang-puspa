package validator

import (
	"backend-golang/internal/adapters/http/dto"
)

type RegistrationValidator interface {
	ValidateRegisterRequest(req *dto.RegistrationRequest) error
}

type registrationValidator struct{}

func NewRegistrationValidator() RegistrationValidator {
	return &registrationValidator{}
}

func (v *registrationValidator) ValidateRegisterRequest(req *dto.RegistrationRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
