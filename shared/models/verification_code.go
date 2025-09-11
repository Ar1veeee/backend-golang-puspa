package models

import (
	"time"
)

type VerificationCode struct {
	Id        int       `json:"id" gorm:"primary_key;auto_increment;"`
	UserId    string    `json:"user_id" gorm:"type:text;not null;index"`
	Code      string    `json:"code" gorm:"type:(200);not null"`
	Status    string    `json:"status" gorm:"type:enum('pending', 'used', 'revoked');default:'pending';not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
