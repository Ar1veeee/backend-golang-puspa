package models

type ObservationAnswer struct {
	Id            int    `json:"id" gorm:"primary_key;type:integer;auto_increment"`
	ObservationId int `json:"observation_id" gorm:"type:integer;not null;index"`
	AspectIndex   string `json:"aspect_index" gorm:"type:varchar(20);not null"`
	Answer        string `json:"answer" gorm:"type:enum('Ya','Tidak');not null"`
	Note          string `json:"note" gorm:"type:text;null"`

	Observation Observation `gorm:"foreignKey:ObservationId;constraint:OnDelete:CASCADE;"`
}
