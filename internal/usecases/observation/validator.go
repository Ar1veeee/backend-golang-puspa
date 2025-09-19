package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/validator"
)

type Validator interface {
	ValidateUpdateScheduledDateRequest(req *dto.UpdateObservationDateRequest) error
}

type observationValidator struct{}

func NewObservationValidator() Validator {
	return &observationValidator{}
}

func (v *observationValidator) ValidateUpdateScheduledDateRequest(req *dto.UpdateObservationDateRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}
