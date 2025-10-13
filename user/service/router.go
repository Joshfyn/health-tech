package service

import "github.com/gin-gonic/gin"

func NewUserRouter(svc *UserService) *gin.Engine {
	r := gin.Default()

	r.POST("/health-tech/users", svc.CreateNewUserNoPwd)
	//r.POST("/health-tech/users/login", svc.SignInUser)
	r.GET("/health-tech/users", svc.GetUsers) // Temp endpoint to get all users
	r.GET("/health-tech/confirm", svc.ConfirmUser)
	

	return r
}
