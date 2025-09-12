package errors

import (
	sharedErrors "backend-golang/shared/errors"
)

var (
	ErrObservationCreationFailed  = sharedErrors.InternalServer("observation_creation_failed", "Gagal membuat akun terapis")
	ErrObservationUpdateFailed    = sharedErrors.InternalServer("observation_update_failed", "failed to update observation")
	ErrObservationDeletionFailed  = sharedErrors.InternalServer("observation_delete_failed", "failed to delete observation")
	ErrObservationRetrievalFailed = sharedErrors.InternalServer("observation_retrieval_failed", "failed to retrieve observations")
)
