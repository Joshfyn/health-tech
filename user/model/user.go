package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRole represents the different user roles in the system
type UserRole string

const (
	RolePatient       UserRole = "patient"
	RoleDoctor        UserRole = "doctor"
	RolePharmacist    UserRole = "pharmacist"
	RoleHospitalAdmin UserRole = "hospital_admin"
	RolePharmacyAdmin UserRole = "pharmacy_admin"
	RoleSuperAdmin    UserRole = "super_admin"
)

// TODO: SA can onboard every one. RHA snd RPA
// TODO: Super Admin can onboard every one. RHA snd RPA

// EntityType represents the type of entity a user is affiliated with
type EntityType string

const (
	EntityHospital EntityType = "Hospital"
	EntityPharmacy EntityType = "Pharmacy"
)

// Contact represents the contact information for a user
type Contact struct {
	Phone   string  `bson:"phone" json:"phone" validate:"required"`
	Address Address `bson:"address" json:"address" validate:"required"`
}

// Address represents a physical address
type Address struct {
	Street  string `bson:"street" json:"street" validate:"required"`
	City    string `bson:"city" json:"city" validate:"required"`
	State   string `bson:"state" json:"state" validate:"required"`
	ZipCode string `bson:"zipCode" json:"zipCode" validate:"required"`
	Country string `bson:"country" json:"country" validate:"required"`
}

// User represents a user in the health-tech system
type User struct {
	ID                uuid.UUID  `bson:"id,omitempty" json:"id,omitempty"`
	Email             string     `bson:"email" json:"email" validate:"required,email"`
	PasswordHash      string     `bson:"passwordHash" json:"-"`
	Role              UserRole   `bson:"role" json:"role" validate:"required,oneof=patient doctor pharmacist hospital_admin pharmacy_admin super_admin"`
	FirstName         string     `bson:"firstName" json:"firstName" validate:"required,min=1,max=50"`
	LastName          string     `bson:"lastName" json:"lastName" validate:"required,min=1,max=50"`
	DateOfBirth       *time.Time `bson:"dateOfBirth,omitempty" json:"dateOfBirth,omitempty"`
	Contact           Contact    `bson:"contact" json:"contact" validate:"required"`
	IsActive          bool       `bson:"isActive" json:"isActive"`
	ConfirmationToken string     `bson:"confirmationToken,omitempty" json:"confirmationToken,omitempty"`
	ConfirmedAt       *time.Time `bson:"confirmedAt,omitempty" json:"confirmedAt,omitempty"`
	LastLogin         *time.Time `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"`
	CreatedAt         time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time  `bson:"updatedAt" json:"updatedAt"`
}

func (u *User) CreateNewUserNoPwd(db *mongo.Database) error {
	collection := db.Collection("users")
	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUsers(db *mongo.Database) ([]User, error) {
	collection := db.Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []User
	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
