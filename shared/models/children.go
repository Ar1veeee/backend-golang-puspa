package models

import (
	"backend-golang/shared/helpers"
	"time"
)

type Children struct {
	Id                 string           `json:"id" gorm:"primaryKey;type:char(26);"`
	ParentId           string           `json:"parent_id" gorm:"type:char(26);not null;index"`
	ChildName          string           `json:"child_name" gorm:"type:varchar(100);not null"`
	ChildGender        string           `json:"child_gender" gorm:"type:enum('Laki-laki', 'Perempuan');not null"`
	ChildBirthPlace    string           `json:"child_birth_place" gorm:"type:varchar(100);not null"`
	ChildBirthDate     helpers.DateOnly `json:"child_birth_date" gorm:"type:date;not null"`
	ChildAddress       []byte           `json:"child_address" gorm:"type:varbinary(500);not null"`
	ChildComplaint     string           `json:"child_complaint" gorm:"type:text;not null"`
	ChildSchool        *string          `json:"child_school" gorm:"type:varchar(100);null"`
	ChildServiceChoice string           `json:"child_service_choice" gorm:"type:varchar(250);not null"`
	ChildReligion      *string          `json:"child_religion" gorm:"type:enum('Islam','Kristen','Katolik','Hindu','Budha','Konghucu','Lainnya');null"`
	CreatedAt          time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"autoUpdateTime"`

	Parent      *Parent      `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
	Observation *Observation `gorm:"foreignKey:ChildId;constraint:OnDelete:CASCADE;"`
}
