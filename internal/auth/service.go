package auth

import (
	error_responses "my-fiber-app/pkg/utils/responses"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Authenticator interface {
	LogIn(userName string, lastName string, browser string, clientIP string) (*AuthLoginResponse, *error_responses.ErrorResponse)
	CheckSession(login_session string, userID float64) (bool, *error_responses.ErrorResponse)
}

type AuthService struct {
	db_pool  *sqlx.DB
	authRepo AuthRepo
	redis    *redis.Client
}

func NewAuthService(db_pool *sqlx.DB, redis *redis.Client) *AuthService {
	r := NewAuthRepoImpl(db_pool, redis)
	return &AuthService{
		db_pool:  db_pool,
		authRepo: r,
	}
}

// Login function
func (a *AuthService) LogIn(userName string, password string, browser string, clientIP string) (*AuthLoginResponse, *error_responses.ErrorResponse) {
	success, errMsg := a.authRepo.LogIn(userName, password, browser, clientIP)
	return success, errMsg
}

// CheckSession function
func (a *AuthService) CheckSession(login_session string, userID float64) (bool, *error_responses.ErrorResponse) {
	success, err := a.authRepo.CheckSession(login_session, userID)
	return success, err
}
