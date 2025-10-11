package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"backend/config"
)

var DB *sql.DB

func InitDB() {
	cfg := config.GetConfig()
	var err error

	DB, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		fmt.Println("DSN used:", cfg.GetDSN())
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
		fmt.Println("DSN used:", cfg.GetDSN())
	}

	log.Println("Database connection successfully established!")
}
