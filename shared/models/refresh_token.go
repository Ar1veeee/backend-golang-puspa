package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	Id        string    `json:"id" gorm:"primary_key;type:char(36);"`
	UserId    string    `json:"userId" gorm:"column:userId;type:char(36);not null;index"`
	Token     string    `json:"token" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expiresAt;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt;autoCreateTime"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`

	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.Id == "" {
		r.Id = uuid.New().String()
	}
	return
}
