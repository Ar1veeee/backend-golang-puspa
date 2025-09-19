package models

import "time"

type ParentDetail struct {
	Id                    string    `gorm:"primary_key;type:char(26);"`
	ParentId              string    `gorm:"type:char(26);not null;index"`
	ParentType            string    `gorm:"type:enum('Ayah','Ibu','Wali');not null;index"`
	ParentName            string    `gorm:"type:varchar(100);not null"`
	ParentPhone           []byte    `gorm:"type:varbinary(100);not null"`
	ParentBirthDate       *string   `gorm:"type:int;null"`
	ParentOccupation      *string   `gorm:"type:varchar(100);null"`
	RelationshipWithChild *string   `gorm:"type:varchar(100);null"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`

	Parent *Parent `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
}
