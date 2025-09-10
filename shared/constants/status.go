package constants

type RegistrationStatus string
type ObservationStatus string
type VerificationCodeStatus string

const (
	RegistrationStatusPending   RegistrationStatus = "pending"
	RegistrationStatusScheduled RegistrationStatus = "scheduled"
	RegistrationStatusCompleted RegistrationStatus = "completed"

	ObservationStatusPending   ObservationStatus = "pending"
	ObservationStatusCompleted ObservationStatus = "completed"

	VerificationCodeStatusPending VerificationCodeStatus = "pending"
	VerificationCodeStatusUsed    VerificationCodeStatus = "used"
	VerificationCodeStatusRevoked VerificationCodeStatus = "revoked"
)
