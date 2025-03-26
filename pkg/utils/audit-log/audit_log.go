package audit_log

import (
	"my-fiber-app/pkg/utils/logs"
	"my-fiber-app/pkg/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func AddUserAuditLog(user_id int, audit_context string, audit_desc string, audit_type_id int, user_agent string, user_name string, ip string, by_id int, db_pool *sqlx.DB) (*bool, error) {

	orderSeqName := "tbl_users_audit_id_seq"
	orderVal, err := sql.GetSeqNextVal(orderSeqName, db_pool)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next order value: %w", err)
	}

	var query = `INSERT INTO tbl_users_audit (
				id, user_id, user_audit_context, user_audit_desc, audit_type_id, user_agent, operator, ip, status_id, "order", created_by, created_at
				) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)`

	app_timezone := os.Getenv("TIME_ZONE")
	location, err := time.LoadLocation(app_timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}
	local_now := time.Now().In(location)
	_, err = db_pool.Exec(
		query,
		*orderVal,
		user_id,
		audit_context,
		audit_desc,
		audit_type_id,
		user_agent,
		user_name,
		ip,
		1,
		*orderVal,
		by_id,
		local_now,
	)
	if err != nil {
		logs.NewCustomLog("user_create_failed", err.Error(), "error")
		return nil, err
	}
	state := true
	return &state, nil
}
