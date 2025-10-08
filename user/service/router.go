package service

import "github.com/gin-gonic/gin"

func NewUserRouter(svc *UserService) *gin.Engine {
	r := gin.Default()

	r.POST("/health-tech/users", svc.CreateNewUserNoPwd)
	r.GET("/health-tech/users", svc.GetUsers) // Temp endpoint to get all users

	return r
}
