package service

import (
	"context"
	"fmt"
	"health-tech/lib/mail"
	"health-tech/user/dto"

	//"health-tech/user/middleware"
	"health-tech/user/model"
	"net/http"

	token "health-tech/lib/tokenmaker"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (svc *UserService) CreateNewUserNoPwd(c *gin.Context) {
	//req := c.MustGet(middleware.ReqBodyCreateNewUserNoPwd).(dto.CreateNewUserNoPwdRequest)
	req := dto.CreateNewUserNoPwdRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		svc.Logger.Errorf(c.Request, "Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate verification token
	token, err := token.GenerateToken(token.TokenRequirements{Email: req.Email})
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error generating verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// generate password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error generating password hash: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var user = model.User{
		ID: func() uuid.UUID {
			id, _ := uuid.NewV4()
			return id
		}(),
		Email:             req.Email,
		Role:              model.UserRole(req.Role),
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DateOfBirth:       req.DateOfBirth,
		Contact:           req.Contact,
		ConfirmationToken: token,
		IsConfirmed:       false,
		PasswordHash:      string(passwordHash),
	}

	session, err := svc.DB.Client().StartSession()
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error starting session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (any, error) {
		err = user.CreateNewUserNoPwd(svc.DB, sessCtx)
		if err != nil {
			svc.Logger.Errorf(c.Request, "Error creating new user: %v", err)
			return nil, fmt.Errorf("create new user failed: %w", err)
		}

		err = mail.SendConfirmationEmail(user.Email, token)
		if err != nil {
			svc.Logger.Errorf(c.Request, "Error sending confirmation email: %v", err)
			return nil, fmt.Errorf("send confirmation email failed: %w", err)
		}

		return nil, nil
	})
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error registering user in transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.CreateNewUserNoPwdResponse{Message: "User created successfully"})
}

func (svc *UserService) GetUsers(c *gin.Context) {
	u := model.User{}
	users, err := u.GetUsers(svc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (svc *UserService) ConfirmUser(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}
	claims, err := token.ValidateVerificationToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := model.User{
		Email: claims.Email,
	}
	err = user.ConfirmUser(svc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User confirmed successfully"})
}
