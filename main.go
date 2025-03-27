package main

import (
	"fmt"
	configs "my-fiber-app/config"
	database "my-fiber-app/config/database"
	redis "my-fiber-app/config/redis"
	"my-fiber-app/handler"
	custom_log "my-fiber-app/pkg/custom_log"
	translate "my-fiber-app/pkg/utils/translate"
	routers "my-fiber-app/routers"
)

func main() {

	// Initial configuration
	app_configs := configs.NewConfig()

	// Initial database
	db_pool := database.GetDB()

	// Initialize router
	app := routers.New(db_pool)

	// Initialize redis client
	rdb := redis.NewRedisClient()

	// Initialize the translate
	if err := translate.Init(); err != nil {
		custom_log.NewCustomLog("Failed_initialize_i18n", err.Err.Error(), "error")
	}

	handler.NewFrontService(app, db_pool, rdb)

	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
}
