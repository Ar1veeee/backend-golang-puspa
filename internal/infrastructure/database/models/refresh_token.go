package models

import (
	"time"
)

type RefreshToken struct {
	Id        int       `gorm:"primary_key;auto_increment;"`
	UserId    string    `gorm:"type:char(36);not null;index"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Revoked   bool      `gorm:"default:false"`

	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
