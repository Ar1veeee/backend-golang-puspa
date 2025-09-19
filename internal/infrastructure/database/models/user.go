package models

import (
	"time"
)

type User struct {
	Id        string    `gorm:"primary_key;type:char(26);"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(200);not null"`
	Role      string    `gorm:"type:enum('Admin', 'Terapis', 'User');default:'User';not null"`
	IsActive  bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	LastLogin time.Time `gorm:"autoUpdateTime;null"`

	RefreshToken []RefreshToken `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Parent       *Parent        `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Therapist    *Therapist     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
