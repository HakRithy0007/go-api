package sql

import (
	"my-fiber-app/pkg/utils/logs"
	"my-fiber-app/pkg/types"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SeqResult struct {
	ID int `db:"id"`
}

func GetUserIdByField(tableName, fieldName string, value interface{}, db *sqlx.DB) (*int, error) {
	// Ensure table and field names are sanitized
	query := fmt.Sprintf(`SELECT id FROM %s WHERE %s = $1 AND deleted_at IS NULL LIMIT 1`, tableName, fieldName)

	var userID *int
	err := db.Get(&userID, query, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}
	return userID, nil
}

func IsExits(tbl_name string, field_name string, value interface{}, db *sqlx.DB) (bool, error) {
	var exists int

	query := fmt.Sprintf(`SELECT 1 as id FROM %s WHERE %s=$1 AND deleted_at IS NULL`, tbl_name, field_name)

	err := db.Get(&exists, query, value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GetSeqNextVal(seqName string, db *sqlx.DB) (*int, error) {
	var result SeqResult
	sql := `SELECT nextval($1) AS id`

	err := db.Get(&result, sql, seqName)
	if err != nil {
		logs.NewCustomLog("failed_to_get_sequence", err.Error(), "error")
		return nil, fmt.Errorf("failed to get sequence value: %w", err)
	}
	return &result.ID, nil
}

func SetSeqNextVal(seqName string, value int, db *sqlx.DB) (*int, error) {

	var result SeqResult

	// Define the SQL query
	sql := `SELECT setval($1, $2) AS id`
	err := db.Get(&result, sql, seqName, value)
	if err != nil {
		logs.NewCustomLog("failed_to_get_sequence_value", err.Error(), "error")
		return nil, fmt.Errorf("failed to set sequence value: %w", err)
	}

	return &result.ID, nil
}

func BuildSQLSort(sorts []share.Sort) string {
	if len(sorts) == 0 {
		return " ORDER BY id"
	}
	var orderClauses []string
	for _, sort := range sorts {
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sort.Property, sort.Direction))
	}
	return " ORDER BY " + strings.Join(orderClauses, ", ")
}

func BuildSQLFilter(req []share.Filter) (string, []interface{}) {
	var sqlFilters []string
	var param []interface{}

	for i, filter := range req {
		paramPlaceholder := fmt.Sprintf("$%d", i+1)

		// Convert the filter value to the appropriate type
		switch v := filter.Value.(type) {
		case string:
			if intValue, err := strconv.Atoi(v); err == nil {
				filter.Value = intValue
			} else if boolValue, err := strconv.ParseBool(v); err == nil {
				filter.Value = boolValue
			} else if dateValue, err := time.Parse("2006-01-02", v); err == nil {
				filter.Value = dateValue
			} else {
				filter.Value = v
			}
		}

		// Handle the converted value
		switch v := filter.Value.(type) {
		case int:
			sqlFilters = append(sqlFilters, fmt.Sprintf("%s = %s", filter.Property, paramPlaceholder))
			param = append(param, v)
		case bool:
			sqlFilters = append(sqlFilters, fmt.Sprintf("%s = %s", filter.Property, paramPlaceholder))
			param = append(param, v)
		case string:
			if strings.Contains(v, "%") {
				// Handle cases with LIKE for wildcard searches
				sqlFilters = append(sqlFilters, fmt.Sprintf("%s LIKE %s", filter.Property, paramPlaceholder))
			} else {
				sqlFilters = append(sqlFilters, fmt.Sprintf("%s = %s", filter.Property, paramPlaceholder))
			}
			param = append(param, v)
		case time.Time:
			// Handle date comparison
			sqlFilters = append(sqlFilters, fmt.Sprintf("%s::DATE = %s", filter.Property, paramPlaceholder))
			param = append(param, v)
		default:
			return "", nil
		}
	}
	filterClause := strings.Join(sqlFilters, " AND ")
	return filterClause, param
}
