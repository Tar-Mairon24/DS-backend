package config

import (
	"fmt"
	"os"
)

type Config struct {
	User     string
	Password string
	Net      string
	Addr     string
	DBName   string
}

func GetConfig() *Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	return &Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Net:      "tcp",
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Net, c.Addr, c.DBName)
}
