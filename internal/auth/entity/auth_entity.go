package entity

import (
	"errors"
	"time"
)

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

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}
	if u.Username == "" {
		return errors.New("username cannot be empty")
	}
	if u.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

type VerificationToken struct {
	Id        int
	UserId    string
	Token     string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (vc *VerificationToken) IsExpired() bool {
	return time.Now().After(vc.ExpiresAt)
}

type RefreshToken struct {
	Id        int
	UserId    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
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

type Parent struct {
	Id        string
	UserId    *string
	TempEmail string
	CreatedAt time.Time
	UpdatedAt time.Time
}
