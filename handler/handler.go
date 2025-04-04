package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	auth "my-fiber-app/internal/auth"
	middleware "my-fiber-app/pkg/middleware"
)

type ServiceHandlers struct {
	Fronted *FrontService
}

type FrontService struct {
	AuthHandler *auth.AuthRoute
}

func NewFrontService(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *FrontService {

	// Auth
	auth := auth.NewAuthRoute(app, db_pool, redis).RegisterAuthRoute()

	// Middleware
	middleware.NewJwtMinddleWare(app, db_pool, redis)

	return &FrontService{
		AuthHandler: auth,
	}
}

func NewServiceHandlers(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *ServiceHandlers {

	return &ServiceHandlers{
		Fronted: NewFrontService(app, db_pool, redis),
	}
}
