package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DATABASE_URL string
}

func GetConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	config := Config{
		PORT: "8080",
	}
	if PORT, exists := os.LookupEnv("PORT"); exists {
		config.PORT = PORT
	}
	if DATABASE_URL, exists := os.LookupEnv("DATABASE_URL"); exists {
		config.DATABASE_URL = DATABASE_URL
	}

	return config
}
