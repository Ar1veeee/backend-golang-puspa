package models

import "time"

type Therapist struct {
	Id               string    `gorm:"primary_key;type:char(26);"`
	UserId           string    `gorm:"type:char(26);null;uniqueIndex"`
	TherapistName    string    `gorm:"type:varchar(100);not null"`
	TherapistSection string    `gorm:"type:enum('Okupasi', 'Fisio', 'Wicara', 'Paedagog');not null"`
	TherapistPhone   []byte    `gorm:"type:varbinary(100);not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	User        *User         `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Observation []Observation `gorm:"foreignKey:TherapistId;constraint:OnDelete:CASCADE;"`
}
