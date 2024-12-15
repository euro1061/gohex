package configs

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
