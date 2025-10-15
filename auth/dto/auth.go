package dto

type SignInUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInUserResponse struct {
	Token string `json:"token"`
}

type ValidateTokenRequest struct {
	Token string `json:"-" validate:"required"`
}

type ValidateTokenResponse struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
