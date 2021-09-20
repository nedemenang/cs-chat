package common

import (
	"os"

	"log"

	"github.com/joho/godotenv"
)

//GetEnv is used to set environment variables
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

//LoadEnvVars loads env variables from .env file
func LoadEnvVars() {
	env := GetEnv("GIN_ENV", "development")

	if env == "production" || env == "staging" {
		log.Println("Not using .env file in production or staging.")
		return
	}

	filename := ".env." + env

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = ".env"
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Println(".env file not loaded")
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
