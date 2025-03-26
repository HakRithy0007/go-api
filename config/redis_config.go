package config

import (
	"log"
	"os"
	utils_env "my-fiber-app/pkg/utils/env"
	"github.com/joho/godotenv"
)

type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisExpire   int
}

func InitRedis() *RedisConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_db := utils_env.GetenvInt("REDIS_DB", 0)
	redis_exprie := utils_env.GetenvInt("REDIS_EXPIRE", 60)

	return &RedisConfig{
		RedisHost:     redis_host,
		RedisPort:     redis_port,
		RedisPassword: redis_password,
		RedisDB:       redis_db,
		RedisExpire:   redis_exprie,
	}
}
