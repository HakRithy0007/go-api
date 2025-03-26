package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	utils_env "my-fiber-app/pkg/utils/env"
)

type AppConfig struct {
	AppHost string
	AppPort int
}

func NewConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	host := os.Getenv("API_HOST")
	port := utils_env.GetenvInt("API_PORT", 8888)

	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}

}
