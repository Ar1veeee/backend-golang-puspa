package dto

import (
	"backend-golang/internal/helpers"
)

type RegistrationRequest struct {
	ChildName          string           `json:"child_name" validate:"required,min=3,max=100"`
	ChildGender        string           `json:"child_gender" validate:"required,oneof=Laki-laki Perempuan"`
	ChildBirthPlace    string           `json:"child_birth_place" validate:"required,min=3,max=100"`
	ChildBirthDate     helpers.DateOnly `json:"child_birth_date" validate:"required" time_format:"2006-01-02"`
	ChildSchool        *string          `json:"child_school"`
	ChildAddress       string           `json:"child_address" validate:"required,min=3,max=100"`
	ChildComplaint     string           `json:"child_complaint" validate:"required,min=3,max=100"`
	ChildServiceChoice string           `json:"child_service_choice" validate:"required"`
	Email              string           `json:"email" validate:"required,email"`
	ParentName         string           `json:"parent_name" validate:"required,min=3,max=100"`
	ParentPhone        string           `json:"parent_phone" validate:"required,min=3,max=100"`
	ParentType         string           `json:"parent_type" validate:"required,oneof=Ayah Ibu Wali"`
}
