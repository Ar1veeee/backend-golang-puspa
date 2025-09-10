package models

import "time"

type Observation struct {
	Id             int       `json:"id" gorm:"primary_key;type:integer;auto_increment;"`
	ChildId        string    `json:"child_id" gorm:"type:char(26);not null;index"`
	TherapistId    string    `json:"therapist_id" gorm:"type:char(26);not null;index"`
	ScheduledDate  string    `json:"scheduled_date" gorm:"type:date;not null"`
	AgeCategory    string    `json:"age_category" gorm:"type:enum('Balita', 'Anak-anak', 'Remaja', 'Lainya');not null"`
	TotalScore     *int      `json:"total_score" gorm:"type:integer;null"`
	Conclusion     *string   `json:"conclusion" gorm:"type:text;null"`
	Recommendation *string   `json:"recommendation" gorm:"type:text;null"`
	Status         string    `json:"status" gorm:"type:enum('Pending','Complete');default:'Pending';not null;index"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Children          Children            `gorm:"foreignKey:ChildId;constraint:OnDelete:CASCADE;"`
	ObservationAnswer []ObservationAnswer `gorm:"foreignKey:ObservationId;constraint:OnDelete:CASCADE;"`
}
