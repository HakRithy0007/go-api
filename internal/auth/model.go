package auth

import (
	utils "my-fiber-app/pkg/utils/helpers"

	"github.com/gofiber/fiber/v2"
)

type AuthLoginRequest struct {
	Auth struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"auth"`
}

type AuthLoginResponse struct {
	Auth struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type"`
	} `json:"auth"`
}

func (r *AuthLoginRequest) bind(c *fiber.Ctx, v *utils.Validator) error {

	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

// own
type UserData struct {
	ID       int    `db:"id"`
	Username string `db:"user_name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	RoleID 	int `db:"role_id"`
}

type RedisSession struct {
	LoginSession string `json:"login_session"`
}
