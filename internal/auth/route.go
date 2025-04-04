package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthRoute struct {
	app     *fiber.App
	handler *AuthHandler
}

func NewAuthRoute(app *fiber.App, dbPool *sqlx.DB, redisClient *redis.Client) *AuthRoute {
	return &AuthRoute{
		app:     app,
		handler: NewAuthHandler(dbPool, redisClient),
	}
}

func (a *AuthRoute) RegisterAuthRoute() *AuthRoute {
	v1 := a.app.Group("/api/v1")
	auth := v1.Group("/auth")
	auth.Post("/login", a.handler.Login)

	return a
}

// POST /auth/register → Register a new user

// POST /auth/login → Login and get a token

// POST /auth/logout → Logout user

// GET /auth/profile → Get user profile

// PUT /auth/profile → Update user profile

// POST /auth/forgot-password → Request password reset

// POST /auth/reset-password → Reset password