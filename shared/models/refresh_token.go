package models

import (
	"time"
)

type RefreshToken struct {
	Id        int       `json:"id" gorm:"primary_key;auto_increment;"`
	UserId    string    `json:"user_id" gorm:"type:char(36);not null;index"`
	Token     string    `json:"token" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`

	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
