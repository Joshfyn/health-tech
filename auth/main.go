package main

import (
	"health-tech/auth/service"
	"health-tech/lib"
	"health-tech/lib/db"
	clog "health-tech/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	authService := service.AuthService{
		DB:     db.ConnectToDB(),
		Logger: *clog.MustNew(service.ServiceName, ""),
	}

	gin.SetMode(gin.ReleaseMode)
	panic(service.NewAuthRouter(&authService).Run(lib.GetServiceHostWithPort(service.ServiceName)))
}
