package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type ResendTokenRequest struct {
	Email string `json:"email" validate:"required,email"`
}
type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3,max=50"`
	Password   string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LoginResponse struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenType    string `json:"tokenType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresAt    string `json:"expiresIn"`
}
