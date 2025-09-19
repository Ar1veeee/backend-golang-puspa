package registration

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/validator"
)

type Validator interface {
	ValidateRegisterRequest(req *dto.RegistrationRequest) error
}

type registrationValidator struct{}

func NewRegistrationValidator() Validator {
	return &registrationValidator{}
}

func (v *registrationValidator) ValidateRegisterRequest(req *dto.RegistrationRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
