package dto

import (
	"health-tech/user/model"
	"time"

	"github.com/gofrs/uuid"
)

type CreateNewUserNoPwdRequest struct {
	UserID       uuid.UUID      `json:"-"`
	PasswordHash string         `json:"-"`
	Email        string         `json:"email" validate:"required,email"`
	Role         model.UserRole `json:"role" validate:"required,oneof=patient super_admin"`
	FirstName    string         `json:"firstName" validate:"required,min=1,max=50"`
	LastName     string         `json:"lastName" validate:"required,min=1,max=50"`
	DateOfBirth  *time.Time     `json:"dateOfBirth" validate:"required"`
	Contact      model.Contact  `json:"contact" validate:"required"`
	NationalID   string         `json:"nationalID" validate:"required"`
	Password     string         `json:"password" validate:"required"`
}

type CreateNewUserNoPwdResponse struct {
	Message string `json:"message"`
}

type CreateProviderAccountRequest struct {
	Email         string         `json:"email"`
	Role          model.UserRole `json:"role"`
	Password      string         `json:"password"`
	UserID        uuid.UUID      `json:"userId"`
	PasswordHash  string         `json:"-"`
	RoleAccountID uuid.UUID      `json:"-"`
}

type PutNewPasswordResponse struct {
	Message string `json:"message"`
}

type ConfirmUserRequest struct {
	Token string `json:"-"`
}

type ConfirmUserResponse struct {
	Message string `json:"message"`
}

type CreateProviderAccountResponse struct {
	Message string `json:"message"`
}
