package entities

import (
	"time"
)

type Admin struct {
	Id         string
	UserId     string
	AdminName  string
	AdminPhone []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time

	User *User
}
