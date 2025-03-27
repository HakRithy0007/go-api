package configs

import (
	"log"
	"os"
	env "my-fiber-app/pkg/utils/env"

	"github.com/joho/godotenv"
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
	host := os.Getenv("API_HOSt")
	port := env.GetenvInt("API_PORT", 8889)
	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}
}
