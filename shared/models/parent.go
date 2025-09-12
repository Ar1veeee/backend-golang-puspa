package models

import "time"

type Parent struct {
	Id                 string    `json:"id" gorm:"primary_key;type:char(26);"`
	UserId             *string   `json:"user_id" gorm:"type:char(26);null;uniqueIndex"`
	TempEmail          string    `json:"temp_email" gorm:"type:varchar(200);not null;uniqueIndex"`
	RegistrationStatus string    `json:"registration_status" gorm:"type:enum('Pending', 'Completed')default:'Pending';not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User         *User          `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	ParentDetail []ParentDetail `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
	Children     []Children     `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
}
