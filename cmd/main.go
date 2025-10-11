package main

import (
	"backend/internal/database"
	"backend/internal/router"
)

func main() {

	database.InitDB()

	ginRouter := router.SetupRouter()

	ginRouter.Run(":8080")
}
