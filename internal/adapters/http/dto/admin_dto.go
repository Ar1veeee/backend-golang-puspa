package dto

type AdminCreateRequest struct {
	Username   string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	AdminName  string `json:"admin_name" validate:"required,min=3,max=100"`
	AdminPhone string `json:"admin_phone" validate:"required,min=3,max=100"`
}

type AdminUpdateRequest struct {
	Username   string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Email      string `json:"email" validate:"omitempty,email"`
	Password   string `json:"password" validate:"omitempty,min=8"`
	AdminName  string `json:"admin_name" validate:"omitempty,min=3,max=100"`
	AdminPhone string `json:"admin_phone" validate:"omitempty,min=3,max=100"`
}

type AdminResponse struct {
	AdminId    string `json:"admin_id" `
	AdminName  string `json:"admin_name" `
	Username   string `json:"username" `
	Email      string `json:"email" `
	AdminPhone string `json:"admin_phone" `
	CreatedAt  string `json:"created_at" `
	UpdatedAt  string `json:"updated_at" `
}
