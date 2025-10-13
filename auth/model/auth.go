package model

import (
	"context"
	"errors"
	"health-tech/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth struct {
	UserID       string         `json:"-"`
	Email        string         `json:"-"`
	Role         model.UserRole `json:"-"`
	PasswordHash string         `json:"-"`
}

func (a *Auth) GetUserByEmail(db *mongo.Database) error {
	var user model.User
	collection := db.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{
		"email":       a.Email,
		"isConfirmed": true,
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	a.UserID = user.ID.String()
	a.PasswordHash = user.PasswordHash
	a.Role = user.Role
	return nil
}
