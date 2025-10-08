package main

import (
	"health-tech/lib"
	"health-tech/lib/db"
	clog "health-tech/logger"
	"health-tech/user/service"

	"github.com/gin-gonic/gin"
)

func main() {

	userService := service.UserService{
		DB:     db.ConnectToDB(),
		Logger: *clog.MustNew(service.ServiceName, ""),
	}

	gin.SetMode(gin.ReleaseMode)
	panic(service.NewUserRouter(&userService).Run(lib.GetServiceHostWithPort(service.ServiceName)))
}
