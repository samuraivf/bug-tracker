package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func OpenPostgres(config *PostgresConfig) (*sql.DB, error) {
	dataSourseName := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open(driver, dataSourseName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
