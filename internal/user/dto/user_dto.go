package dto

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,max=15"`
}
type UserUpdateRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty" validate:"omitempty,max=15"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	CreateAt string `json:"createdAt"`
	UpdateAt string `json:"updatedAt"`
}
