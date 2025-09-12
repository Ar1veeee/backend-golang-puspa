package entity

import (
	"backend-golang/shared/helpers"
	"time"
)

type Parent struct {
	Id                 string
	UserId             *string
	TempEmail          string
	RegistrationStatus string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ParentDetail struct {
	Id          string
	ParentId    string
	ParentType  string
	ParentName  string
	ParentPhone []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
}

type Observation struct {
	Id            int
	ChildId       string
	ScheduledDate helpers.DateOnly
	AgeCategory   string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
