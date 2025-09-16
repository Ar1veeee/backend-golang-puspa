package entities

import (
	"errors"
	"time"
)

type RefreshToken struct {
	Id        int
	UserId    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() error {
	if rt.Revoked {
		return errors.New("refresh token has been revoked")
	}
	if rt.IsExpired() {
		return errors.New("refresh token has expired")
	}
	return nil
}
