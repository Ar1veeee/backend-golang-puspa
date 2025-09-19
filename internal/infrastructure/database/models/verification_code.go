package models

import (
	"time"
)

type VerificationCode struct {
	Id        int       `gorm:"primary_key;auto_increment;"`
	UserId    string    `gorm:"type:text;not null;index"`
	Code      string    `gorm:"type:(200);not null"`
	Status    string    `gorm:"type:enum('Pending', 'Used', 'Revoked');default:'Pending';not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
