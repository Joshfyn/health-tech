package service

import (
	"health-tech/user/middleware"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(svc *UserService) *gin.Engine {
	r := gin.Default()

	r.POST("/health-tech/user", middleware.ValidateCreateNewUserAccountRequest, svc.CreateNewUserAccount)
	r.GET("/health-tech/user", svc.GetUsers)
	r.GET("/health-tech/user/confirm", middleware.ValidateConfirmUserRequest, svc.ConfirmUser)
	r.POST("/health-tech/user/healthcare-providers", middleware.ValidateCreateProviderAccountRequest, svc.CreateProviderAccount)

	return r
}
