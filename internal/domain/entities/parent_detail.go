package entities

import "time"

type ParentDetail struct {
	Id          string
	ParentId    string
	ParentType  string
	ParentName  string
	ParentPhone []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Parent *Parent
}
