package usecase

import (
	"backend-golang/internal/registration/delivery/http/dto"
	"backend-golang/shared/helpers"
	"fmt"
)

type RegistrationValidator interface {
	validateRegisterRequest(req *dto.RegistrationRequest) error
}

type registrationValidator struct{}

func NewRegistrationValidator() RegistrationValidator {
	return &registrationValidator{}
}

func (v *registrationValidator) validateRegisterRequest(req *dto.RegistrationRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.ChildAge < 0 {
		return fmt.Errorf("child age cannot be negative")
	}

	return nil
}
