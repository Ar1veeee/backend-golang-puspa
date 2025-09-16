package entities

import (
	"time"
)

type Therapist struct {
	Id               string
	UserId           string
	TherapistName    string
	TherapistSection string
	TherapistPhone   []byte
	CreatedAt        time.Time
	UpdatedAt        time.Time

	User *User
}
