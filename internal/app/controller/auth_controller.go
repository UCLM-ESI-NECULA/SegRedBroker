package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seg-red-broker/internal/app/client"
	"seg-red-broker/internal/app/common"
	"seg-red-broker/internal/app/dao"
	"seg-red-broker/internal/app/service"
)

type AuthControllerImpl struct {
	svc service.AuthService
}

func NewAuthController() *AuthControllerImpl {
	return &AuthControllerImpl{svc: service.NewAuthService(*client.NewAuthClient())}
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	ValidateToken(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (ac *AuthControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ac.Signup)
	router.POST("/login", ac.Login)
	router.POST("/checkToken", ac.ValidateToken)
}

// Signup handles the /signup endpoint
func (ac *AuthControllerImpl) Signup(c *gin.Context) {
	user, err := checkUserInput(c)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, err.Error())
		return
	}
	token, err := ac.svc.Signup(user.Username, user.Password)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, "invalid credentials")
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// Login handles the /login endpoint
func (ac *AuthControllerImpl) Login(c *gin.Context) {
	user, err := checkUserInput(c)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, err.Error())
		return
	}
	token, err := ac.svc.Login(user.Username, user.Password)
	if err != nil {
		common.NewAPIError(c, http.StatusUnauthorized, err, "invalid credentials")
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func (ac *AuthControllerImpl) ValidateToken(c *gin.Context) {
	token := checkTokenInput(c)
	if token == "" {
		return
	}
	username, err := ac.svc.ValidateToken(token)
	if err != nil {
		common.NewAPIError(c, http.StatusUnauthorized, err, "invalid token")
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})

}

// checkUserInput checks if the user input is valid
func checkUserInput(c *gin.Context) (dao.User, error) {
	var user dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		common.NewAPIError(c, http.StatusUnauthorized, err, "error when mapping request")
		return user, err
	}
	if user.Username == "" || user.Password == "" {
		return user, fmt.Errorf("username and password are required")
	}
	return user, nil
}

func checkTokenInput(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		common.NewAPIError(c, http.StatusUnauthorized, errors.New("unauthorized"), "authorization header is required")
		return ""
	}
	return token
}
