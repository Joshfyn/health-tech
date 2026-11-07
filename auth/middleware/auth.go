package middleware

import (
	"health-tech/auth/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateTokenRequest is the request body for the ValidateToken endpoint
// Token is sent to the endpoint in the Authorization header as `X-USER-TOKEN`
func ValidateTokenRequest(c *gin.Context) {
	token := c.GetHeader("X-USER-TOKEN")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	req := dto.ValidateTokenRequest{
		Token: token,
	}

	c.Set(ReqBodyValidateToken, req)
	c.Next()
}

func ValidateSignInUserRequest(c *gin.Context) {
	req := dto.SignInUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Email == "" {
		log.Printf("Email is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if req.Password == "" {
		log.Printf("Password is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	c.Set(ReqBodySignInUser, req)
	c.Next()
}
