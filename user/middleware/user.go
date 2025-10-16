package middleware

import (
	"health-tech/user/dto"
	"health-tech/user/model"
	"log"
	"net/http"

	token "health-tech/lib/tokenmaker"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SuperAdminOnly(c *gin.Context) {
	authToken := c.GetHeader("X-USER-TOKEN")
	if authToken == "" {
		log.Printf("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	authUser, err := token.ValidateVerificationToken(authToken)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if authUser.Role != string(model.RoleSuperAdmin) {
		log.Printf("Only super admins can create provider accounts")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	c.Set(ReqBodyAuthUser, authUser)
	c.Next()
}

func AdminsOnly(c *gin.Context) {
	authToken := c.GetHeader("X-USER-TOKEN")
	if authToken == "" {
		log.Printf("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	authUser, err := token.ValidateVerificationToken(authToken)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if authUser.Role != string(model.RoleHospitalAdmin) && authUser.Role != string(model.RolePharmacyAdmin) && authUser.Role != string(model.RoleSuperAdmin) {
		log.Printf("Only admins can create provider accounts")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	c.Set(ReqBodyAuthUser, authUser)
	c.Next()
}

// ValidateTokenRequest is the request body for the ValidateToken endpoint
// Token is sent to the endpoint in the Authorization header as `X-USER-TOKEN`
func ValidateCreateProviderAccountRequest(c *gin.Context) {
	req := dto.CreateProviderAccountRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Email == "" {
		log.Printf("Email is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if req.UserID == uuid.Nil {
		log.Printf("User ID is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	if req.Role != model.RoleDoctor && req.Role != model.RolePharmacist && req.Role != model.RoleHospitalAdmin && req.Role != model.RolePharmacyAdmin {
		log.Printf("Invalid role: %s Only doctors, pharmacists, hospital admins, and pharmacy admins can be created", req.Role)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if req.Password == "" {
		log.Printf("Password is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	req.PasswordHash = string(passwordHash)
	roleAccountID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating user ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	req.RoleAccountID = roleAccountID
	c.Set(ReqBodyCreateProviderAccount, req)
	c.Next()
}

func ValidateConfirmUserRequest(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	req := dto.ConfirmUserRequest{
		Token: token,
	}
	c.Set(ReqBodyConfirmUser, req)
	c.Next()
}

func ValidateCreateNewUserAccountRequest(c *gin.Context) {
	req := dto.CreateNewUserNoPwdRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Password == "" {
		log.Printf("Password is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	req.PasswordHash = string(passwordHash)

	userID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating user ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	req.UserID = userID
	c.Set(ReqBodyCreateNewUserAccount, req)
	c.Next()
}
