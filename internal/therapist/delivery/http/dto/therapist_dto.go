package dto

type TherapistCreateRequest struct {
	Username         string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	Role             string `json:"role" validate:"required,max=15"`
	TherapistName    string `json:"therapist_name" validate:"required,min=3,max=100"`
	TherapistSection string `json:"therapist_section" validate:"required,oneof=Okupasi Fisio Wicara Paedagog"`
	TherapistPhone   string `json:"therapist_phone" validate:"required,min=3,max=100"`
}
type TherapistUpdateRequest struct {
	Username         string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email            string `json:"email,omitempty" validate:"omitempty,email"`
	Password         string `json:"password,omitempty"`
	Role             string `json:"role,omitempty" validate:"omitempty,max=15"`
	TherapistName    string `json:"therapist_name,omitempty" validate:"omitempty,min=3,max=100"`
	TherapistSection string `json:"therapist_section,omitempty" validate:"omitempty,oneof=Okupasi Fisio Wicara Paedagog"`
	TherapistPhone   string `json:"therapist_phone,omitempty" validate:"omitempty,min=3,max=100"`
}

type TherapistResponse struct {
    Id               string `json:"id" db:"users.id"`              
    Username         string `json:"username" db:"users.username"`
    Email            string `json:"email" db:"users.email"`
    Role             string `json:"role" db:"users.role"`
    TherapistId      string `json:"therapist_id" db:"therapists.id"` 
    TherapistName    string `json:"therapist_name" db:"therapists.therapist_name"`
    TherapistSection string `json:"therapist_section" db:"therapists.therapist_section"`
    TherapistPhone   string `json:"therapist_phone" db:"therapists.therapist_phone"`
    CreatedAt        string `json:"created_at" db:"therapists.created_at"`
    UpdatedAt        string `json:"updated_at" db:"therapists.updated_at"`
}
