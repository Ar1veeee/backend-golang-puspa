package errors

import (
	sharedErrors "backend-golang/shared/errors"
)

var (
	ErrRegistrationFailed = sharedErrors.InternalServer("registration_failed", "failed to create registration")
)
