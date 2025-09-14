package main

import (
	"backend/internal/database"
	"backend/internal/router"
)

func main() {
	// Initialize the database connection
	database.InitDB()

	ginRouter := router.SetupRouter()

	// Start the server
	ginRouter.Run(":8080")
}
