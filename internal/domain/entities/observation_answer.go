package entities

type ObservationAnswer struct {
	Id            int
	ObservationId int
	QuestionId    int
	Answer        bool
	ScoreEarned   int
	Note          *string

	Observation         *Observation
	ObservationQuestion *ObservationQuestion
}
