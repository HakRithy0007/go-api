package auth

import (
	http_status_code "my-fiber-app/pkg/http_status_codes"
	responses "my-fiber-app/pkg/utils/responses"
	translate "my-fiber-app/pkg/translate"
	utils "my-fiber-app/pkg/utils/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// AuthHandler struct
type AuthHandler struct {
	db_pool     *sqlx.DB
	authService Authenticator
	redis       *redis.Client
}

func NewHandler(db_pool *sqlx.DB, redis *redis.Client) *AuthHandler {
	return &AuthHandler{
		db_pool:     db_pool,
		authService: NewAuthService(db_pool, redis),
	}
}

// Login function
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	v := utils.NewValidator()
	req := &AuthLoginRequest{}

	// Bind and validate
	if err := req.bind(c, v); err != nil {
		msg, err_msg := translate.TranslateWithError(c, "login_invalid")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				responses.NewResponseError(
					err_msg.ErrorString(),
					http_status_code.Translate_failed,
					err_msg.Err,
				),
			)
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewResponseError(
				msg,
				http_status_code.Login_invalid,
				err,
			),
		)
	}

	browser := c.Get("User-Agent", "unknown")
	clientIP := c.IP()
	success, err := h.authService.LogIn(req.Auth.Username, req.Auth.Password, browser, clientIP)
	if err != nil {
		msg, msg_err := translate.TranslateWithError(c, err.MessageID)
		if msg_err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(responses.NewResponseError(
				msg_err.Err.Error(),
				http_status_code.Translate_failed,
				msg_err.Err,
			))
		}
		return c.Status(fiber.StatusUnauthorized).JSON(responses.NewResponseError(
			msg,
			http_status_code.Login_failed,
			err.Err,
		))
	} else {
		msg, err_msg := translate.TranslateWithError(c, "login_success")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(responses.NewResponseError(
				err_msg.ErrorString(),
				http_status_code.Translate_failed,
				err_msg.Err,
			))
		}
		return c.Status(fiber.StatusOK).JSON(responses.NewResponse(
			msg,
			http_status_code.Login_success,
			success,
		))
	}
}
