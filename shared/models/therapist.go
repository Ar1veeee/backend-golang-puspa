package models

import "time"

type Therapist struct {
	Id               string    `json:"id" gorm:"primary_key;type:char(26);"`
	UserId           string    `json:"user_id" gorm:"type:char(26);null;uniqueIndex"`
	TherapistName    string    `json:"therapist_name" gorm:"type:varchar(100);not null"`
	TherapistSection string    `json:"therapist_section" gorm:"type:enum('Okupasi', 'Fisio', 'Wicara', 'Paedagog');not null"`
	TherapistPhone   []byte    `json:"therapist_phone" gorm:"type:varbinary(100);not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User        *User         `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Observation []Observation `gorm:"foreignKey:TherapistId;constraint:OnDelete:CASCADE;"`
}
