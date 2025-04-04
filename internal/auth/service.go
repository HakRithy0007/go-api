package auth

import (
	"my-fiber-app/pkg/utils/error"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Login(username, password string) (*AuthResponse, *error.ErrorResponse)
	CheckSession(loginSession string, userID float64) (bool, *error.ErrorResponse)
}

type authServiceImpl struct {
	repo AuthRepository
}

func NewAuthService(dbPool *sqlx.DB, redisClient *redis.Client) AuthService {
	repo := NewAuthRepository(dbPool, redisClient)
	return &authServiceImpl{
		repo: repo,
	}
}

// Login
func (a *authServiceImpl) Login(username, password string) (*AuthResponse, *error.ErrorResponse) {
	return a.repo.Login(username, password)
}

func (a *authServiceImpl) CheckSession(loginSession string, userID float64) (bool, *error.ErrorResponse) {
	return a.repo.CheckSession(loginSession, userID)
}