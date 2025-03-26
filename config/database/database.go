package database

import (
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	once    sync.Once
	db_pool *sqlx.DB
)

func InitializeDB() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	var err_db error

	db_pool, err_db = sqlx.Connect("postgres", DATABASE_URL)
	if err_db != nil {
		log.Fatalln("Error connection to the database", err_db)
	}

	// Set connection pool settings
	db_pool.SetMaxIdleConns(10)
	db_pool.SetMaxOpenConns(10)
	db_pool.SetConnMaxLifetime(0)
}

func GetDB() *sqlx.DB {
	once.Do(InitializeDB)
	return db_pool
}
