package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"seg-red-broker/internal/app/config"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file: ", err)
		panic(err)
	}
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")
	app := config.SetupRouter()

	err := app.Run(":" + port)
	if err != nil {
		log.Error("Error running server: ", err)
	}
}
