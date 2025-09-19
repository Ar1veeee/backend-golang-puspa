package entities

import "time"

type ObservationQuestion struct {
	Id             int
	QuestionCode   string
	AgeCategory    string
	QuestionNumber int
	QuestionText   string
	Score          int
	IsActive       bool
	CreatedAt      time.Time

	ObservationAnswer []*ObservationAnswer
}
