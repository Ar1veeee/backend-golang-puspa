package dto

import "backend-golang/shared/helpers"

type PendingObservationsResponse struct {
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

type CompleteObservationsResponse struct {
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

type ObservationDetailResponse struct {
	ObservationId  int              `json:"observation_id"`
	AgeCategory    string           `json:"age_category"`
	ChildName      string           `json:"child_name"`
	ChildBirthDate helpers.DateOnly `json:"child_birth_date"`
	ChildSchool    string           `json:"child_school"`
	ChildAge       int              `json:"child_age"`
	ChildAddress   string           `json:"child_address"`
	ChildComplaint string           `json:"child_complaint"`
	ParentName     string           `json:"parent_name"`
	ParentPhone    string           `json:"parent_phone"`
	Status         string           `json:"status"`
}
