package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seg-red-broker/internal/app/service"
)

type BrokerControllerImpl struct {
	svc service.BrokerService
}

func NewBrokerController(r *gin.RouterGroup) *BrokerControllerImpl {
	c := &BrokerControllerImpl{svc: service.NewBrokerService()}
	c.RegisterRoutes(r)
	return c
}

type BrokerController interface {
	GetVersion(c *gin.Context)
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

// RegisterRoutes registers the broker routes
func (ac *BrokerControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/version", ac.GetVersion)
}

// GetVersion handles the /version endpoint
func (ac *BrokerControllerImpl) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, ac.svc.GetVersion())
}
