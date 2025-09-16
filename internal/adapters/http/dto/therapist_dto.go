package dto

type TherapistCreateRequest struct {
	Username         string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	TherapistName    string `json:"therapist_name" validate:"required,min=3,max=100"`
	TherapistSection string `json:"therapist_section" validate:"required,oneof=Okupasi Fisio Wicara Paedagog"`
	TherapistPhone   string `json:"therapist_phone" validate:"required,min=3,max=100"`
}
type TherapistUpdateRequest struct {
	Username         string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email            string `json:"email,omitempty" validate:"omitempty,email"`
	Password         string `json:"password,omitempty"`
	TherapistName    string `json:"therapist_name,omitempty" validate:"omitempty,min=3,max=100"`
	TherapistSection string `json:"therapist_section,omitempty" validate:"omitempty,oneof=Okupasi Fisio Wicara Paedagog"`
	TherapistPhone   string `json:"therapist_phone,omitempty" validate:"omitempty,min=3,max=100"`
}

type TherapistResponse struct {
	UserId           string `json:"user_id"`
	TherapistId      string `json:"therapist_id"`
	TherapistName    string `json:"therapist_name"`
	TherapistSection string `json:"therapist_section"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	TherapistPhone   string `json:"therapist_phone"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
