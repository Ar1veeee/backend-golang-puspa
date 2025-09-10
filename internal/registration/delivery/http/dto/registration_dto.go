package dto

import (
	"encoding/json"
	"time"
)

type RegistrationRequest struct {
	ChildName          string    `json:"child_name" validate:"required,min=3,max=100"`
	ChildGender        string    `json:"child_gender" validate:"required,oneof=Laki-laki Perempuan"`
	ChildBirthPlace    string    `json:"child_birth_place" validate:"required,min=3,max=100"`
	ChildBirthDate     time.Time `json:"child_birth_date" validate:"required" time_format:"2006-01-02"`
	ChildAge           int       `json:"child_age" validate:"required"`
	ChildSchool        string    `json:"child_school"`
	ChildAddress       string    `json:"child_address" validate:"required,min=3,max=100"`
	ChildComplaint     string    `json:"child_complaint" validate:"required,min=3,max=100"`
	ChildServiceChoice string    `json:"child_service_choice" validate:"required"`
	Email              string    `json:"email" validate:"required,email"`
	ParentName         string    `json:"parent_name" validate:"required,min=3,max=100"`
	ParentPhone        string    `json:"parent_phone" validate:"required,min=3,max=100"`
	ParentType         string    `json:"parent_type" validate:"required,oneof=Ayah Ibu Wali"`
}

func (r *RegistrationRequest) UnmarshalJSON(b []byte) error {
	type Alias RegistrationRequest
	aux := &struct {
		*Alias
		ChildBirthDate string `json:"child_birth_date"`
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	if aux.ChildBirthDate != "" {
		t, err := time.Parse("2006-01-02", aux.ChildBirthDate)
		if err != nil {
			return err
		}
		r.ChildBirthDate = t
	}
	return nil
}
