package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment")
	}
	return os.Getenv(key)
}
