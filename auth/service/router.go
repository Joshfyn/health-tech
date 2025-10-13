package service

import "github.com/gin-gonic/gin"

func NewAuthRouter(svc *AuthService) *gin.Engine {
	r := gin.Default()

	r.POST("/health-tech/auth/signin", svc.SignInUser)

	//TODO: r.GET("/health-tech/auth/validate", svc.ValidateToken)

	//TODO: r.POST("/health-tech/auth/refresh", svc.RefreshToken)

	return r
}
