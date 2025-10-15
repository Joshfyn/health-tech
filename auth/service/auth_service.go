package service

import (
	clog "health-tech/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ServiceName = "auth"
)

type AuthService struct {
	DB     *mongo.Database
	Logger clog.Clog
}
