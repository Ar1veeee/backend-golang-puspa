package constants

type Role string
type RegistrationStatus string
type ObservationStatus string
type VerificationCodeStatus string

const (
	RoleAdmin     Role = "Admin"
	RoleUser      Role = "User"
	RoleTherapist Role = "Terapis"
)

const (
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusComplete RegistrationStatus = "Complete"

	ObservationStatusPending   ObservationStatus = "Pending"
	ObservationStatusCompleted ObservationStatus = "Complete"

	VerificationCodeStatusPending VerificationCodeStatus = "Pending"
	VerificationCodeStatusUsed    VerificationCodeStatus = "Used"
	VerificationCodeStatusRevoked VerificationCodeStatus = "Revoked"
)
