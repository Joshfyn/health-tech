package model

import (
	"context"
	"errors"
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
	ID           uuid.UUID  `bson:"id,omitempty" json:"id,omitempty"`
	PrimaryEmail string     `bson:"email" json:"email" validate:"required,email"` // set unique index
	FirstName    string     `bson:"firstName" json:"firstName" validate:"required,min=1,max=50"`
	LastName     string     `bson:"lastName" json:"lastName" validate:"required,min=1,max=50"`
	DateOfBirth  *time.Time `bson:"dateOfBirth,omitempty" json:"dateOfBirth,omitempty"`
	Contact      Contact    `bson:"contact" json:"contact" validate:"required"`
	NationalID   string     `bson:"nationalID" json:"nationalID" validate:"required"`
	CreatedAt    time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `bson:"updatedAt" json:"updatedAt"`
}

type RoleAccount struct {
	RoleAccountID     uuid.UUID `bson:"roleAccountID" json:"roleAccountID" validate:"required"`
	UserID            uuid.UUID `bson:"userId" json:"userId" validate:"required"`
	Email             string    `bson:"email" json:"email" validate:"required,email"`
	Role              UserRole  `bson:"role" json:"role" validate:"required,oneof=patient doctor pharmacist hospital_admin pharmacy_admin super_admin"`
	IsActive          bool      `bson:"isActive" json:"isActive" validate:"required"`
	IsConfirmed       bool      `bson:"isConfirmed" json:"isConfirmed" validate:"required"`
	ConfirmationToken string    `bson:"confirmationToken" json:"confirmationToken" validate:"required"`
	PasswordHash      string    `bson:"passwordHash" json:"passwordHash" validate:"required"`
	LastLogin         time.Time `bson:"lastLogin" json:"lastLogin" validate:"required"`
	CreatedAt         time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (u *User) CreateNewUser(db *mongo.Database, sessCtx mongo.SessionContext) error {
	collection := db.Collection("user")
	_, err := collection.InsertOne(sessCtx, u)
	if err != nil {
		return err
	}

	return nil
}

func (u *RoleAccount) CreateNewUserAccount(db *mongo.Database, sessCtx mongo.SessionContext) error {
	collection := db.Collection("roleAccount")
	_, err := collection.InsertOne(sessCtx, u)
	if err != nil {
		return err
	}
	return nil
}

func (u *RoleAccount) CreateNewUserAccountWithoutTransaction(db *mongo.Database) error {
	collection := db.Collection("roleAccount")
	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUsers(db *mongo.Database) ([]User, error) {
	collection := db.Collection("user")
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

func (u *User) GetSingleUser(db *mongo.Database) error {
	collection := db.Collection("user")
	var user User
	if u.PrimaryEmail == "" && u.ID == uuid.Nil {
		return errors.New("email and/or id must be provided")
	}
	filter := bson.M{}
	if u.PrimaryEmail != "" {
		filter = bson.M{"email": u.PrimaryEmail}
	}
	if u.ID != uuid.Nil {
		filter["id"] = u.ID
	}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *RoleAccount) ConfirmUser(db *mongo.Database) error {
	collection := db.Collection("roleAccount")
	filter := bson.M{"email": u.Email, "userId": u.UserID, "role": u.Role}
	update := bson.M{"$set": bson.M{"isConfirmed": true}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *RoleAccount) GetAllUserRoles(db *mongo.Database) (map[UserRole]any, error) {
	collection := db.Collection("roleAccount")
	filter := bson.M{"userId": u.UserID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	userRoles := make(map[UserRole]any)
	for cursor.Next(context.Background()) {
		var userAccount RoleAccount
		err := cursor.Decode(&userAccount)
		if err != nil {
			return nil, err
		}
		userRoles[userAccount.Role] = nil
	}
	return userRoles, nil
}
