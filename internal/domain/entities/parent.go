package entities

import "time"

type Parent struct {
	Id                 string
	UserId             *string
	TempEmail          string
	RegistrationStatus string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	User         *User
	ParentDetail []ParentDetail
	Children     []Children
}
