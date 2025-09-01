package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        string    `json:"id" gorm:"primary_key;type:char(36);"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt;autoUpdateTime"`

	RefreshToken []RefreshToken `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}
