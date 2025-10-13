package dto

import (
	"health-tech/user/model"
	"time"
)

type CreateNewUserNoPwdRequest struct {
	Email       string        `json:"email" validate:"required,email"`
	Role        string        `json:"role" validate:"required,oneof=patient doctor pharmacist hospital_admin pharmacy_admin super_admin"`
	FirstName   string        `json:"firstName" validate:"required,min=1,max=50"`
	LastName    string        `json:"lastName" validate:"required,min=1,max=50"`
	DateOfBirth *time.Time    `json:"dateOfBirth" validate:"required"`
	Contact     model.Contact `json:"contact" validate:"required"`
	Password    string        `json:"password" validate:"required"`
}

type CreateNewUserNoPwdResponse struct {
	Message string `json:"message"`
}
