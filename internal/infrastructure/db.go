package infrastructure

import (
	"database/sql"
	"fmt"
	"user-management/config"
	"user-management/internal/gateways"

	mysql "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func (s MySQLStore) Query(q string, args ...any) (*sql.Rows, error) {
	return s.db.Query(q, args...)
}

func getConfig() *mysql.Config {
	env := config.GetConfig()
	config := mysql.NewConfig()

	config.User = env.MySQLUsername
	config.Passwd = env.MySQLPassword
	config.Net = "tcp"
	config.Addr = env.DBContainer
	config.DBName = env.MySQLDatabase

	return config
}

func NewMySQLStore() (gateways.Store, error) {
	db, err := sql.Open("mysql", getConfig().FormatDSN())

	if err != nil {
		return MySQLStore{}, fmt.Errorf("failed to connect db: %v", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	return MySQLStore{db: db}, nil
}
