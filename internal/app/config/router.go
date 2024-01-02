package config

import (
	"github.com/gin-gonic/gin"
	"seg-red-broker/internal/app/common"
	"seg-red-broker/internal/app/controller"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(common.GlobalErrorHandler())

	v1 := r.Group("/api/v1")

	brokerCtrl := controller.NewBrokerController()
	brokerCtrl.RegisterRoutes(v1)

	fileCtrl := controller.NewFileController()
	fileCtrl.RegisterRoutes(v1)

	authCtrl := controller.NewAuthController()
	authCtrl.RegisterRoutes(v1)
	return r
}
