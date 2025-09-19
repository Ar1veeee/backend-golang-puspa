package models

type ObservationAnswer struct {
	Id            int     `gorm:"primary_key;type:integer;auto_increment"`
	ObservationId int     `gorm:"type:integer;not null;index"`
	QuestionId    int     `gorm:"type:integer;not null"`
	Answer        bool    `gorm:"type:bool;not null"`
	ScoreEarned   int     `gorm:"type:integer;not null;default:0"`
	Note          *string `gorm:"type:text;null"`

	Observation         Observation         `gorm:"foreignKey:ObservationId;constraint:OnDelete:CASCADE;"`
	ObservationQuestion ObservationQuestion `gorm:"foreignKey:QuestionId;constraint:OnDelete:CASCADE;"`
}
