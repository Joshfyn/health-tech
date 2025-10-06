package services

import (
	clog "health-tech/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	DB     *mongo.Database
	Logger clog.Clog
}

func (u *UserService) NewUserService(db *mongo.Database) {
	u.Logger.Errorf(nil, "NewUserService")
}
