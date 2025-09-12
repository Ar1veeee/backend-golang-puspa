package models

import (
	"time"
)

type User struct {
	Id        string    `json:"id" gorm:"primary_key;type:char(26);"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"type:varchar(200);not null"`
	Role      string    `json:"role" gorm:"type:enum('Admin', 'Terapis', 'User');default:'User';not null"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	LastLogin time.Time `json:"last_login" gorm:"autoUpdateTime;null"`

	RefreshToken []RefreshToken `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Parent       *Parent        `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Therapist    *Therapist     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
