package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var SqlDB *sql.DB

func GetDB() *sql.DB {
	logrus.Info("Getting SQL DB connection")
	if SqlDB == nil {
		logrus.Error("SQL DB connection is not initialized")
	}
	return SqlDB
}

func Init() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	logrus.Infof("Connecting to database at %s:%s/%s", dbHost, dbPort, dbName)

	// Open SQL connection
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	// Test the connection
	err = database.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to ping database")
	}

	SqlDB = database
	logrus.Info("Successfully connected to database")
}

func GetDBErrorNoRows(err error) bool {
	return err == sql.ErrNoRows
}
