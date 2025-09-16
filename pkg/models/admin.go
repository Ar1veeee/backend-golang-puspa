package models

import "time"

type Admin struct {
	Id         string    `gorm:"primary_key;type:char(26);"`
	UserId     string    `gorm:"type:char(26);null;uniqueIndex"`
	AdminName  string    `gorm:"type:varchar(100);not null"`
	AdminPhone []byte    `gorm:"type:varbinary(100);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	User *User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
