package dto

type SignInUserRequest struct {
	Email       string        `json:"email" validate:"required,email"`
	Password    string        `json:"password" validate:"required"`
}

type SignInUserResponse struct {
	Token string `json:"token"`
}
