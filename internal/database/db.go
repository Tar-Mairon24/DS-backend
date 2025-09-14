package database

import (
	"backend/config"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := config.GetConfig()
	var err error

	DB, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	log.Println("Database connection successfully established!")
}
