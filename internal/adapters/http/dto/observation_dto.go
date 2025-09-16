package dto

import (
	"backend-golang/internal/helpers"
)

type ObservationsResponse struct {
	ObservationId  int              `json:"observation_id"`
	AgeCategory    string           `json:"age_category"`
	ChildName      string           `json:"child_name"`
	ChildAge       int              `json:"child_age"`
	ChildSchool    *string          `json:"child_school"`
	ChildComplaint string           `json:"child_complaint"`
	ParentName     string           `json:"parent_name"`
	ParentPhone    string           `json:"parent_phone"`
	ScheduledDate  helpers.DateOnly `json:"scheduled_date"`
	Status         string           `json:"status"`
}

type DetailObservationResponse struct {
	ObservationId  int              `json:"observation_id"`
	ChildName      string           `json:"child_name"`
	ChildBirthDate helpers.DateOnly `json:"child_birth_date"`
	ChildAge       int              `json:"child_age"`
	ChildGender    string           `json:"child_gender"`
	ChildSchool    *string          `json:"child_school"`
	ChildAddress   string           `json:"child_address"`

	ParentName  string `json:"parent_name"`
	ParentType  string `json:"parent_type"`
	ParentPhone string `json:"parent_phone"`
	Email       string `json:"email"`

	ChildComplaint string `json:"child_complaint"`
}

type DetailObservationRequest struct {
	ObservationId int `json:"observation_id" validate:"required"`
}
