package models

import (
	"backend-golang/internal/helpers"
	"time"
)

type Children struct {
	Id                 string           `gorm:"primaryKey;type:char(26);"`
	ParentId           string           `gorm:"type:char(26);not null;index"`
	ChildName          string           `gorm:"type:varchar(100);not null"`
	ChildGender        string           `gorm:"type:enum('Laki-laki', 'Perempuan');not null"`
	ChildBirthPlace    string           `gorm:"type:varchar(100);not null"`
	ChildBirthDate     helpers.DateOnly `gorm:"type:date;not null"`
	ChildAddress       []byte           `gorm:"type:varbinary(500);not null"`
	ChildComplaint     string           `gorm:"type:text;not null"`
	ChildSchool        *string          `gorm:"type:varchar(100);null"`
	ChildServiceChoice string           `gorm:"type:varchar(250);not null"`
	ChildReligion      *string          `gorm:"type:enum('Islam','Kristen','Katolik','Hindu','Budha','Konghucu','Lainnya');null"`
	CreatedAt          time.Time        `gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime"`

	Parent      *Parent      `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
	Observation *Observation `gorm:"foreignKey:ChildId;constraint:OnDelete:CASCADE;"`
}
