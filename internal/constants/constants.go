package constants

type ContextKey string

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
	ContextUserID   ContextKey = "userId"
	ContextUserRole ContextKey = "userRole"
)

const (
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusComplete RegistrationStatus = "Complete"

	ObservationStatusPending   ObservationStatus = "Pending"
	ObservationStatusScheduled ObservationStatus = "Scheduled"
	ObservationStatusCompleted ObservationStatus = "Complete"

	VerificationCodeStatusPending VerificationCodeStatus = "Pending"
	VerificationCodeStatusUsed    VerificationCodeStatus = "Used"
	VerificationCodeStatusRevoked VerificationCodeStatus = "Revoked"
)
