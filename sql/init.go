package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func initMigrations(dialect string, db *sql.DB) {
	if err := goose.SetDialect(dialect); err != nil {
		logrus.Fatalf("failed to set %v dialect: %v", dialect, err)
	}
	if version, err := goose.EnsureDBVersion(db); err == nil {
		logrus.Infof("version of %s migrations: %v", dialect, version)
	}
	if err := goose.Up(db, dialect); err != nil {
		logrus.Fatalf("failed to make %v migrations up: %v", dialect, err)
	}
}

//go:embed clickhouse/*.sql
var clickhouseMigrations embed.FS

//go:embed mysql/*.sql
var mysqlMigrations embed.FS

//go:embed postgres/*.sql
var potgresqlMigrations embed.FS

func InitClickHouseMigrations(db *sql.DB) {
	goose.SetBaseFS(clickhouseMigrations)
	goose.SetTableName("goose_db_version")
	initMigrations(string(goose.DialectClickHouse), db)
}

func InitMySQLMigrations(db *sql.DB) {
	goose.SetBaseFS(mysqlMigrations)
	goose.SetTableName("goose_db_version")
	initMigrations(string(goose.DialectMySQL), db)
}

func InitPostgreSQLMigrations(db *sql.DB) {
	goose.SetBaseFS(potgresqlMigrations)
	goose.SetTableName("public.goose_db_version")
	initMigrations(string(goose.DialectPostgres), db)
}

type UnusedTypeToBeDeleted struct{}
