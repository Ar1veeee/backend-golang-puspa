package models

type ObservationAnswer struct {
	Id            int    `gorm:"primary_key;type:integer;auto_increment"`
	ObservationId int    `gorm:"type:integer;not null;index"`
	AspectIndex   string `gorm:"type:varchar(20);not null"`
	Answer        string `gorm:"type:enum('Ya','Tidak');not null"`
	Note          string `gorm:"type:text;null"`

	Observation Observation `gorm:"foreignKey:ObservationId;constraint:OnDelete:CASCADE;"`
}
