package entities

import (
	"backend-golang/internal/helpers"
	"time"
)

type Children struct {
	Id                 string
	ParentId           string
	ChildName          string
	ChildGender        string
	ChildBirthPlace    string
	ChildBirthDate     helpers.DateOnly
	ChildAddress       []byte
	ChildComplaint     string
	ChildSchool        *string
	ChildServiceChoice string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Parent      *Parent
	Observation *Observation
}
