package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	PORT               string
	DATABASE_URL       string
	VALIDATE_JWT       bool
	JWT_SIGNING_SECRET string
	DEBUG_LOGGING      bool
}

var Config = &Configuration{
	PORT:               "8080",
	DATABASE_URL:       "",
	VALIDATE_JWT:       true,
	JWT_SIGNING_SECRET: "",
	DEBUG_LOGGING:      false,
}

func InitConfiguration() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := godotenv.Load(pwd + "/.env"); err != nil {
		log.Print("No .env file found")
	}
	if PORT, exists := os.LookupEnv("PORT"); exists {
		Config.PORT = PORT
	}
	if DATABASE_URL, exists := os.LookupEnv("DATABASE_URL"); exists {
		Config.DATABASE_URL = DATABASE_URL
	}
	if VALIDATE_JWT, exists := os.LookupEnv("VALIDATE_JWT"); exists {
		Config.VALIDATE_JWT = VALIDATE_JWT != "false"
	}
	if JWT_SIGNING_SECRET, exists := os.LookupEnv("JWT_SIGNING_SECRET"); exists {
		Config.JWT_SIGNING_SECRET = JWT_SIGNING_SECRET
	}
	if DEBUG_LOGGING, exists := os.LookupEnv("DEBUG_LOGGING"); exists {
		Config.DEBUG_LOGGING = DEBUG_LOGGING == "true"
	}
}
