package service

import (
	"health-tech/auth/middleware"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(svc *AuthService) *gin.Engine {
	r := gin.Default()

	r.POST("/health-tech/auth/signin", svc.SignInUser)

	r.GET("/health-tech/auth/validate", middleware.ValidateTokenRequest, svc.ValidateToken)

	r.POST("/health-tech/auth/refresh", middleware.ValidateTokenRequest, svc.RefreshToken)

	return r
}
