package entities

import (
	"backend-golang/internal/helpers"
	"time"
)

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

	Children *Children
}
