package entity

import (
	"time"
)

type Children struct {
	Id                 string
	ParentId           string
	ChildName          string
	ChildGender        string
	ChildBirthPlace    string
	ChildBirthDate     string
	ChildAge           int
	ChildAddress       []byte
	ChildComplaint     string
	ChildSchool        *string
	ChildServiceChoice string
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

type Parent struct {
	Id                 string
	TempEmail          string
	RegistrationStatus string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Observation struct {
	Id            int
	ChildId       string
	ScheduledDate string
	AgeCategory   string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
