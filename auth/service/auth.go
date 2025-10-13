package service

import (
	"health-tech/auth/dto"
	"health-tech/auth/model"
	"log"
	"net/http"

	token "health-tech/lib/tokenmaker"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (svc *AuthService) SignInUser(c *gin.Context) {
	// Takes the user email and password
	req := dto.SignInUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		svc.Logger.Errorf(c.Request, "Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify the email exists and the password is correct
	var auth = model.Auth{
		Email:        req.Email,
	}
	log.Println("auth", auth)
	err := auth.GetUserByEmail(svc.DB)
	if err != nil {
		svc.Logger.
			WithField("email", req.Email).
			Errorf(c.Request, "Error getting user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// decrypt the password from the db
	err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(req.Password))
	if err != nil {
		svc.Logger.
		WithField("email", req.Email).
		Errorf(c.Request, "Error comparing password hash: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// generate token
	tokenStr, err := token.GenerateToken(token.TokenRequirements{Email: req.Email, UserID: auth.UserID})
	if err != nil {
		svc.Logger.
			WithField("email", req.Email).
			WithField("user_id", auth.UserID).
			Errorf(c.Request, "Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.SignInUserResponse{Token: tokenStr})
}
