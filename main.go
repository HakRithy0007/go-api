package main

import (
	"fmt"
	configs "my-fiber-app/config"
	database "my-fiber-app/config/database"
	redis "my-fiber-app/config/redis"
	logs "my-fiber-app/pkg/utils/logs"
	translate "my-fiber-app/pkg/translate"
	routes "my-fiber-app/routes"
	handler "my-fiber-app/handler"
)

func main() {

	// Initial configuration
	app_configs := configs.NewConfig()

	// Initial database
	db_pool := database.GetDB()

	// Initialize router
	app := routes.New(db_pool)

	// Initialize redis client
	rdb := redis.NewRedisClient()

	// Initialize the translate
	if err := translate.Init(); err != nil {
		logs.NewCustomLog("Failed_initialize_i18n", err.Err.Error(), "error")
	}

	handler.NewFrontService(app, db_pool, rdb)

	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
}
