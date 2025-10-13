package service

import (
	clog "health-tech/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ServiceName = "user"
)

type UserService struct {
	DB     *mongo.Database
	Logger clog.Clog
}

func (svc *UserService) NewUserService(db *mongo.Database) {
	svc.Logger.Errorf(nil, "NewUserService")
}
