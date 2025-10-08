package service

import (
	"health-tech/user/dto"
	//"health-tech/user/middleware"
	"health-tech/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (svc *UserService) CreateNewUserNoPwd(c *gin.Context) {
	//req := c.MustGet(middleware.ReqBodyCreateNewUserNoPwd).(dto.CreateNewUserNoPwdRequest)
	req := dto.CreateNewUserNoPwdRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		svc.Logger.Errorf(c.Request, "Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user = model.User{
		ID: func() uuid.UUID {
			id, _ := uuid.NewV4()
			return id
		}(),
		Email:       req.Email,
		Role:        model.UserRole(req.Role),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Contact:     req.Contact,
	}

	err := user.CreateNewUserNoPwd(svc.DB)
	if err != nil {
		svc.Logger.Errorf(c.Request, "Error binding JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
