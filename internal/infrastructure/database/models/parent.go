package models

import "time"

type Parent struct {
	Id                 string    `gorm:"primary_key;type:char(26);"`
	UserId             *string   `gorm:"type:char(26);null;uniqueIndex"`
	TempEmail          string    `gorm:"type:varchar(200);not null;uniqueIndex"`
	RegistrationStatus string    `gorm:"type:enum('Pending', 'Completed')default:'Pending';not null"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`

	User         *User          `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	ParentDetail []ParentDetail `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
	Children     []Children     `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
}
