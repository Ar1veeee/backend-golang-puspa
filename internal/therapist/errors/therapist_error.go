package errors

import (
	sharedErrors "backend-golang/shared/errors"
)

var (
	ErrTherapistCreationFailed  = sharedErrors.InternalServer("therapist_creation_failed", "Gagal membuat akun terapis")
	ErrTherapistUpdateFailed    = sharedErrors.InternalServer("therapist_update_failed", "failed to update therapist")
	ErrTherapistDeletionFailed  = sharedErrors.InternalServer("therapist_delete_failed", "failed to delete therapist")
	ErrTherapistRetrievalFailed = sharedErrors.InternalServer("therapist_retrieval_failed", "failed to retrieve therapists")
)
