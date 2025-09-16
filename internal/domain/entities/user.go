package entities

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
	LastLogin time.Time

	VerificationToken []VerificationToken
	RefreshToken      []RefreshToken
	Parent            []Parent
	Admin             []Admin
}
