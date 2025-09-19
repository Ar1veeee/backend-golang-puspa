package dto

import "backend-golang/internal/helpers"

type ChildResponse struct {
	ChildId        string           `json:"child_id"`
	ChildName      string           `json:"child_name"`
	ChildBirthDate helpers.DateOnly `json:"child_birth_date"`
	ChildGender    string           `json:"child_gender"`
	ParentName     string           `json:"parent_name"`
	ParentPhone    string           `json:"parent_phone"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
}
