package configs

import (
	"log"
	"os"
	utils "my-fiber-app/pkg/utils"

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
	port := utils.GetenvInt("API_PORT", 8889)
	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}
}
