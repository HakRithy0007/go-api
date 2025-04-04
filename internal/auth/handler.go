package auth

import (
	constants "my-fiber-app/pkg/constants"
	"my-fiber-app/pkg/utils/response"
	"my-fiber-app/pkg/utils/translate"
	custom_validator "my-fiber-app/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(dbPool *sqlx.DB, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		authService: NewAuthService(dbPool, redisClient),
	}
}

// Login
func (a *AuthHandler) Login(c *fiber.Ctx) error {
	v := custom_validator.NewValidator()
	req := &AuthLoginRequest{}

	if err := req.bind(c, v); err != nil {
		msg, errMsg := translate.TranslateWithError(c, "login_invalid")
		if errMsg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					errMsg.ErrorString(),
					constants.Translate_failed,
					errMsg.Err,
				),
			)
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewResponseError(
				msg,
				constants.Login_invalid,
				err,
			),
		)
	}
	// Service
	success, err := a.authService.Login(req.Auth.Username, req.Auth.Password)

	if err != nil {
		msg, msgErr := translate.TranslateWithError(c, err.MessageID)
		if msgErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewResponseError(
				msgErr.Err.Error(),
				constants.Translate_failed,
				msgErr.Err,
			))
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response.NewResponseError(
			msg,
			constants.LoginFailed,
			err.Err,
		))
	}

	msg, errMsg := translate.TranslateWithError(c, "login_success")
	if errMsg != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewResponseError(
			errMsg.ErrorString(),
			constants.Translate_failed,
			errMsg.Err,
		))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(
		msg,
		constants.Login_success,
		success,
	))
}
