package errors

import "errors"

var (
	ErrTherapistCreationFailed  = errors.New("failed to create therapist")
	ErrTherapistUpdateFailed    = errors.New("failed to update therapist")
	ErrTherapistDeletionFailed  = errors.New("failed to delete therapist")
	ErrTherapistRetrievalFailed = errors.New("failed to retrieve therapists")
)
