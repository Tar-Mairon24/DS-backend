package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"backend/internal/database"
	"backend/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using environment variables or defaults")
	}

	timezone := os.Getenv("SERVER_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
		log.Print("SERVER_TIMEZONE not set, defaulting to UTC")
	}
	loc, err := time.LoadLocation(timezone)
	log.Print("Setting server timezone to ", timezone)
	if err != nil {
		log.Print("Failed to load timezone ", timezone, ", defaulting to UTC")
		loc = time.UTC
	}
	time.Local = loc

	database.InitDB()

	ginRouter := router.SetupRouter()

	ginRouter.Run(":8080")
}
