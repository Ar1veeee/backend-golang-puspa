package models

import (
	"backend-golang/internal/helpers"
	"time"
)

type Observation struct {
	Id             int              `gorm:"primary_key;type:integer;auto_increment;"`
	ChildId        string           `gorm:"type:char(26);not null;index"`
	TherapistId    string           `gorm:"type:char(26);not null;index"`
	ScheduledDate  helpers.DateOnly `gorm:"type:date;not null"`
	AgeCategory    string           `gorm:"type:enum('Balita', 'Anak-anak', 'Remaja', 'Lainnya');not null"`
	TotalScore     int              `gorm:"type:integer;null"`
	Conclusion     string           `gorm:"type:text;null"`
	Recommendation string           `gorm:"type:text;null"`
	Status         string           `gorm:"type:enum('Pending','Complete');default:'Pending';not null;index"`
	CreatedAt      time.Time        `gorm:"autoCreateTime"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime"`

	Children          *Children           `gorm:"foreignKey:ChildId;constraint:OnDelete:CASCADE;"`
	ObservationAnswer []ObservationAnswer `gorm:"foreignKey:ObservationId;constraint:OnDelete:CASCADE;"`
}
