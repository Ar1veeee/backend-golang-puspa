package entities

import "time"

type VerificationToken struct {
	Id        int
	UserId    string
	Token     string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

func (vc *VerificationToken) IsExpired() bool {
	return time.Now().After(vc.ExpiresAt)
}
