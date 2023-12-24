package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed clickhouse/*.sql
var clickhouseMigrations embed.FS

func InitClickHouseMigrations(db *sql.DB) {
	goose.SetBaseFS(clickhouseMigrations)
	if err := goose.SetDialect(string(goose.DialectClickHouse)); err != nil {
		logrus.Fatal("failed to set goose dialect: ", err)
	}
	if err := goose.Up(db, "clickhouse"); err != nil {
		logrus.Fatal("failed to make clickhouse migrations up: ", err)
	}
}

//go:embed mysql/*.sql
var mysqlMigrations embed.FS

func InitMySQLMigrations(db *sql.DB) {
	goose.SetBaseFS(mysqlMigrations)
	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		logrus.Fatal("failed to set goose dialect: ", err)
	}
	if err := goose.Up(db, "mysql"); err != nil {
		logrus.Fatal("failed to make mysql migrations up: ", err)
	}
}
