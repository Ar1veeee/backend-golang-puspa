package constants

type RegistrationStatus string
type ObservationStatus string
type VerificationCodeStatus string

const (
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusComplete RegistrationStatus = "Complete"

	ObservationStatusPending   ObservationStatus = "Pending"
	ObservationStatusCompleted ObservationStatus = "Complete"

	VerificationCodeStatusPending VerificationCodeStatus = "Pending"
	VerificationCodeStatusUsed    VerificationCodeStatus = "Used"
	VerificationCodeStatusRevoked VerificationCodeStatus = "Revoked"
)
