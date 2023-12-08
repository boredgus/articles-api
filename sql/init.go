package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitMigrations(db *sql.DB, dialect string) {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(dialect); err != nil {
		logrus.Fatal("failed to set goose dialect: ", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		logrus.Fatal("failed to make migrations up: ", err)
	}
}
