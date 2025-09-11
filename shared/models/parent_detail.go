package models

import "time"

type ParentDetail struct {
	Id                    string    `json:"id" gorm:"primary_key;type:char(26);"`
	ParentId              string    `json:"parent_id" gorm:"type:char(26);not null;index"`
	ParentType            string    `json:"parent_type" gorm:"type:enum('Ayah','Ibu','Wali');not null;index"`
	ParentName            string    `json:"parent_name" gorm:"type:varchar(100);not null"`
	ParentPhone           []byte    `json:"contact" gorm:"type:varbinary(100);not null"`
	ParentAge             int       `json:"parent_age" gorm:"type:int;null"`
	ParentOccupation      *string   `json:"parent_occupation" gorm:"type:varchar(100);null"`
	RelationshipWithChild string    `json:"relationship_with_child" gorm:"type:varchar(100);null"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Parent Parent `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE;"`
}
