package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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
	certs := os.Getenv("CERTS_FOLDER")
	app := config.SetupRouter()
	err := app.RunTLS(":"+port, filepath.Join(certs, "mycert.crt"), filepath.Join(certs, "mycert.key"))
	if err != nil {
		panic(err)
	}
}
