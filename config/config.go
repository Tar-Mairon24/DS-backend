package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Password string
	Net      string
	Addr     string
	DBName   string
}

func GetConfig() *Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory")
	}

	envPath := filepath.Join(wd, "..", ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	return &Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("USER_PASSWORD"),
		Net:      os.Getenv("tcp"),
		Addr:     os.Getenv("DB_PORT:3606"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Net, c.Addr, c.DBName)
}
