package auth

import (
	logs "my-fiber-app/pkg/utils/logs"
	redis_utils "my-fiber-app/pkg/redis"
	error_responses "my-fiber-app/pkg/utils/responses"
	utils "my-fiber-app/pkg/utils/env"
	audit_log "my-fiber-app/pkg/utils/audit-log"
	"context"
	"database/sql"
	"errors" 
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	LogIn(userName string, password string, userAgent string, clientIP string) (*AuthLoginResponse, *error_responses.ErrorResponse)
	CheckSession(login_session string, userID float64) (bool, *error_responses.ErrorResponse)
}

type AuthRepoImpl struct {
	db_pool *sqlx.DB
	redis   *redis.Client
}

func NewAuthRepoImpl(db *sqlx.DB, redis *redis.Client) *AuthRepoImpl {
	return &AuthRepoImpl{
		db_pool: db,
		redis:   redis,
	}
}

// Login function
func (a *AuthRepoImpl) LogIn(userName string, password string, userAgent string, clientIP string) (*AuthLoginResponse, *error_responses.ErrorResponse) {

	var user UserData
	msg := error_responses.ErrorResponse{}

	query :=
		`
					SELECT 
						id, 
						user_name, 
						email, 
						password, 
						role_id
					FROM tbl_users 
					WHERE user_name = $1 AND password = $2 AND deleted_at is NULL
				`

	err := a.db_pool.Get(&user, query, userName, password)
	if err != nil {
		logs.NewCustomLog("user_not_found", err.Error(), "error")
		return nil, msg.NewErrorResponse("user_not_found", fmt.Errorf("User not found. Please check the provided information."))
	}

	var r AuthLoginResponse

	hours := utils.GetenvInt("JWT_EXP_HOUR", 7)
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)
	login_session, err := uuid.NewV7()
	if err != nil {
		logs.NewCustomLog("uuid_generate_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("uuid_generate_failed", fmt.Errorf("Failed to generate UUID. Please try again later."))
	}

	// Create the JWT claims
	claims := jwt.MapClaims{
		"user_id":       user.ID,
		"username":      user.Username,
		"login_session": login_session.String(),
		"exp":           expirationTime.Unix(),
		"role_id":       user.RoleID,
	}

	// Set redis data
	key := fmt.Sprintf("user_info_id:%d", user.ID)
	rdb := redis_utils.NewRedisUtil(a.redis)
	rdb.SetCacheKey(key, claims, context.Background())

	errs := godotenv.Load()
	if errs != nil {
		logs.NewCustomLog("error_load_env", errs.Error(), "error")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")
	updateQuery := `UPDATE 
						tbl_users 
					SET 
						login_session = $1 WHERE id = $2`

	_, err = a.db_pool.Exec(updateQuery, login_session.String(), user.ID)
	if err != nil {
		logs.NewCustomLog("session_update_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("session_update_failed", fmt.Errorf("cannot update session"))
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		logs.NewCustomLog("jwt_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("jwt_failed", fmt.Errorf("failed to get jwt"))
	}

	r.Auth.Token = tokenString
	r.Auth.TokenType = "jwt"

	// Audit log
	auditDesc := fmt.Sprintf(`User : %s has been login the dashboard`, userName)
	_, err = audit_log.AddUserAuditLog(user.ID, "Login", auditDesc, 1, userAgent, user.Username, clientIP, user.ID, a.db_pool)
	if err != nil {
		logs.NewCustomLog("add_audit_log_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("add_audit_log_failed", fmt.Errorf("cannot insert data to audit log"))
	}

	return &r, nil
}

// CheckSession function
func (a *AuthRepoImpl) CheckSession(login_session string, userID float64) (bool, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	key := fmt.Sprintf("user:%d", int(userID))
	rdb := redis_utils.NewRedisUtil(a.redis)

	key_data, err := rdb.GetCacheKey(key, context.Background())
	if err == nil {
		if key_data.LoginSession == login_session {
			return true, nil
		}
	}

	var dbLoginSession string

	query :=
		`
        SELECT
			login_session
        FROM tbl_users
        WHERE login_session = $1
    `
	err = a.db_pool.Get(&dbLoginSession, query, login_session)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logs.NewCustomLog("invalid_sesion_id", "invalid login session: "+login_session, "warn")
			return false, msg.NewErrorResponse("invalid_sesion_id", fmt.Errorf("invalid login session"))
		}
		logs.NewCustomLog("query_data_failed", err.Error(), "error")
		return false, msg.NewErrorResponse("query_data_failed", fmt.Errorf("database query error"))
	}

	if dbLoginSession != login_session {
		return false, msg.NewErrorResponse("invalid_sesion_id", fmt.Errorf("invalid login session"))
	}

	return true, nil
}
