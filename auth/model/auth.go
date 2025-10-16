package model

import (
	"context"
	"errors"
	"health-tech/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth struct {
	RoleAccountId string         `json:"-"`
	Email         string         `json:"-"`
	Role          model.UserRole `json:"-"`
	PasswordHash  string         `json:"-"`
}

func (a *Auth) GetUserByEmail(db *mongo.Database) error {
	var user model.RoleAccount
	collection := db.Collection("roleAccount")
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
	a.RoleAccountId = user.RoleAccountID.String()
	a.PasswordHash = user.PasswordHash
	a.Role = user.Role
	return nil
}
