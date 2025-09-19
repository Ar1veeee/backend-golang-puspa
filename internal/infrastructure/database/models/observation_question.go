package models

import "time"

type ObservationQuestion struct {
	Id             int       `gorm:"primary_key;type:integer;auto_increment"`
	QuestionCode   string    `gorm:"type:varchar(6);unique;not null"`
	AgeCategory    string    `gorm:"type:enum('Balita', 'Anak-anak', 'Remaja', 'Lainya');not null"`
	QuestionNumber int       `gorm:"type:integer;not null"`
	QuestionText   string    `gorm:"type:text;not null"`
	Score          int       `gorm:"type:integer;not null"`
	IsActive       bool      `gorm:"type:bool;default:true"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	ObservationAnswer []ObservationAnswer `gorm:"foreignKey:QuestionId;constraint:OnDelete:CASCADE;"`
}
