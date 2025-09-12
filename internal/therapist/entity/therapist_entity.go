package entity

import "time"

type User struct {
	Id        string
	Username  string
	Email     string
	Password  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Therapist struct {
	Id               string
	UserId           string
	TherapistName    string
	TherapistSection string
	TherapistPhone   []byte
	CreatedAt        time.Time
	UpdatedAt        time.Time

	User *User `json:"user,omitempty"`
}
