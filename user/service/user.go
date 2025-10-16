package service

import (
	"context"
	"fmt"
	"health-tech/lib/mail"
	"health-tech/user/dto"
	"health-tech/user/middleware"

	//"health-tech/user/middleware"
	"health-tech/user/model"
	"net/http"

	token "health-tech/lib/tokenmaker"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateNewUserAccount creates a new user account for only super admin and patients
func (svc *UserService) CreateNewUserAccount(c *gin.Context) {
	req := c.MustGet(middleware.ReqBodyCreateNewUserAccount).(dto.CreateNewUserNoPwdRequest)
	// generate verification token
	token, err := token.GenerateToken(token.TokenRequirements{Email: req.Email, UserID: req.UserID.String(), Role: string(req.Role)})
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error generating verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var user = model.User{
		ID:           req.UserID,
		PrimaryEmail: req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DateOfBirth:  req.DateOfBirth,
		Contact:      req.Contact,
		NationalID:   req.NationalID,
	}

	session, err := svc.DB.Client().StartSession()
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error starting session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (any, error) {
		err = user.CreateNewUser(svc.DB, sessCtx)
		if err != nil {
			svc.Logger.Errorf(c.Request, "Error creating new user: %v", err)
			return nil, fmt.Errorf("create new user failed: %w", err)
		}

		if req.Role == model.RolePatient || req.Role == model.RoleSuperAdmin {
			// automatically create a patient and super admin account
			patientAccount := model.UserAccount{
				Email:        user.PrimaryEmail,
				Role:         req.Role,
				UserID:       user.ID,
				PasswordHash: req.PasswordHash,
				IsConfirmed:  false,
				IsActive:     true,
			}
			err = patientAccount.CreateNewUserAccount(svc.DB, sessCtx)
			if err != nil {
				svc.Logger.Errorf(c.Request, "Error creating patient account: %v", err)
				return nil, fmt.Errorf("create patient account failed: %w", err)
			}
		}
		err = mail.SendConfirmationEmail(user.PrimaryEmail, token)
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

func (svc *UserService) CreateProviderAccount(c *gin.Context) {
	req := c.MustGet(middleware.ReqBodyCreateProviderAccount).(dto.CreateProviderAccountRequest)

	// verify the email exists in user
	user := model.User{
		PrimaryEmail: req.Email,
	}
	err := user.GetSingleUser(svc.DB)
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error getting user: %v, user must have a basic account first", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user must have a basic account first"})
		return
	}

	token, err := token.GenerateToken(token.TokenRequirements{Email: req.Email, UserID: req.UserAccountID.String(), Role: string(req.Role)})
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error generating verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session, err := svc.DB.Client().StartSession()
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error starting session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (any, error) {
		account := model.UserAccount{
			UserAccountID:     req.UserAccountID,
			Email:             req.Email,
			Role:              req.Role,
			UserID:            user.ID,
			PasswordHash:      req.PasswordHash,
			ConfirmationToken: token,
			IsConfirmed:       false,
			IsActive:          true,
		}
		err = account.CreateNewUserAccount(svc.DB, sessCtx)
		if err != nil {
			svc.Logger.Errorf(c.Request, "Error creating patient account: %v", err)
			return nil, fmt.Errorf("create patient account failed: %w", err)
		}

		err = mail.SendConfirmationEmail(req.Email, token)
		if err != nil {
			svc.Logger.Errorf(c.Request, "Error sending confirmation email: %v", err)
			return nil, fmt.Errorf("send confirmation email failed: %w", err)
		}

		return nil, nil
	})
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error creating provider account in transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.CreateProviderAccountResponse{Message: "Provider account created successfully"})
}

func (svc *UserService) ConfirmUser(c *gin.Context) {
	req := c.MustGet(middleware.ReqBodyConfirmUser).(dto.ConfirmUserRequest)
	claims, err := token.ValidateVerificationToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify the user exists in user
	user := model.User{
		PrimaryEmail: claims.Email,
		ID:           uuid.FromStringOrNil(claims.UserID),
	}
	err = user.GetSingleUser(svc.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	userAccount := model.UserAccount{
		Email:             claims.Email,
		Role:              model.UserRole(claims.Role),
		UserID:            uuid.FromStringOrNil(claims.UserID),
		ConfirmationToken: req.Token,
		IsConfirmed:       true,
	}
	err = userAccount.ConfirmUser(svc.DB)
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error confirming user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ConfirmUserResponse{Message: "User confirmed successfully"})
}
