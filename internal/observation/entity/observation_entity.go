package entity

import (
	"backend-golang/shared/helpers"
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

	Parent *Parent `json:"parent,omitempty"`
}

type Parent struct {
	Id                 string
	UserId             *string
	TempEmail          string
	RegistrationStatus string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	User         *User          `json:"user,omitempty"`
	ParentDetail []ParentDetail `json:"parent_details,omitempty"`
	Children     []Children     `json:"children,omitempty"`
}

type ParentDetail struct {
	Id          string
	ParentId    string
	ParentType  string
	ParentName  string
	ParentPhone []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Parent *Parent `json:"parent,omitempty"`
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

	Parent      *Parent      `json:"parent,omitempty"`
	Observation *Observation `json:"observation,omitempty"`
}

type Observation struct {
	Id             int
	ChildId        string
	TherapistId    string
	AgeCategory    string
	TotalScore     int
	Conclusion     string
	Recommendation string
	ScheduledDate  helpers.DateOnly
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Children *Children `json:"children,omitempty"`
}
