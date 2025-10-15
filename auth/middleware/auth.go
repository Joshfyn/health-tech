package middleware

import (
	"health-tech/auth/dto"
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
