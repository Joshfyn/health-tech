package service

import (
	"health-tech/auth/dto"
	"health-tech/auth/middleware"
	"health-tech/auth/model"
	"net/http"

	token "health-tech/lib/tokenmaker"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (svc *AuthService) SignInUser(c *gin.Context) {
	// Takes the user email and password
	req := c.MustGet(middleware.ReqBodySignInUser).(dto.SignInUserRequest)

	// verify the email exists and the password is correct
	var auth = model.Auth{
		Email: req.Email,
	}

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
	tokenStr, err := token.GenerateToken(token.TokenRequirements{Email: req.Email, UserID: auth.RoleAccountId, Role: string(auth.Role)})
	if err != nil {
		svc.Logger.
			WithField("email", req.Email).
			WithField("user_id", auth.RoleAccountId).
			Errorf(c.Request, "Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.SignInUserResponse{Token: tokenStr})
}

func (svc *AuthService) ValidateToken(c *gin.Context) {
	req := c.MustGet(middleware.ReqBodyValidateToken).(dto.ValidateTokenRequest)
	claims, err := token.ValidateVerificationToken(req.Token)
	if err != nil {
		svc.Logger.
			WithField("token", req.Token).
			Errorf(c.Request, "Error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}
	c.JSON(http.StatusOK, dto.ValidateTokenResponse{Email: claims.Email, Role: claims.Role, UserID: claims.UserID})
}

func (svc *AuthService) RefreshToken(c *gin.Context) {
	req := c.MustGet(middleware.ReqBodyValidateToken).(dto.ValidateTokenRequest)
	claims, err := token.ValidateVerificationToken(req.Token)
	if err != nil {
		svc.Logger.
			WithField("token", req.Token).
			Errorf(c.Request, "Error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenStr, err := token.GenerateToken(token.TokenRequirements{Email: claims.Email, UserID: claims.UserID, Role: claims.Role})
	if err != nil {
		svc.Logger.
			WithField("email", claims.Email).
			WithField("user_id", claims.UserID).
			Errorf(c.Request, "Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{Token: tokenStr})
}
