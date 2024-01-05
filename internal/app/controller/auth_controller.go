package controller

import (
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

func NewAuthController(g *gin.RouterGroup) *AuthControllerImpl {
	controller := &AuthControllerImpl{svc: service.NewAuthService(*client.NewAuthClient())}
	controller.RegisterRoutes(g)
	return controller
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	//ValidateToken(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (ac *AuthControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ac.Signup)
	router.POST("/login", ac.Login)
	//router.POST("/checkToken", ac.ValidateToken)
}

// Signup handles the /signup endpoint
func (ac *AuthControllerImpl) Signup(c *gin.Context) {
	// Check input
	user, err := checkUserInput(c)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Create user
	token, err := ac.svc.Signup(user.Username, user.Password)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, token)
}

// Login handles the /login endpoint
func (ac *AuthControllerImpl) Login(c *gin.Context) {
	// Check input
	user, err := checkUserInput(c)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Login user
	token, err := ac.svc.Login(user.Username, user.Password)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, token)
}

func (ac *AuthControllerImpl) ValidateToken(c *gin.Context) {
	// Check input
	user, err := CheckTokenInput(c, ac.svc)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// checkUserInput checks if the user input is valid
func checkUserInput(c *gin.Context) (*dao.User, error) {
	var user *dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		return nil, common.BadRequestError("invalid request body")
	}
	if user.Username == "" {
		return nil, common.EmptyParamsError("username")
	}
	if user.Password == "" {
		return nil, common.EmptyParamsError("password")
	}
	return user, nil
}

func CheckTokenInput(c *gin.Context, svc service.AuthService) (*dao.User, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, common.UnauthorizedError("authorization header is required")
	}
	user, err := svc.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
