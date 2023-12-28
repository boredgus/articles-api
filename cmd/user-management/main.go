// Articles API
//
// Spec Documentation for article service.
//
// Version: 1.2.0
//
// Schemes:
//   - http
//
// BasePath: /
// Host: localhost:8080
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"database/sql"
	infrastructure "user-management/cmd/user-management/internal"
	"user-management/config"
	"user-management/pkg/db"
	migrations "user-management/sql"
)

func init() {
	config.InitConfig()
}

//go:generate swagger generate spec --scan-models --input=../../docs/init.json --output=../../docs/swagger.json
func main() {
	router := infrastructure.GetRouter(
		infrastructure.NewAppController(
			db.NewMySQLStore(func(db *sql.DB) { migrations.InitMySQLMigrations(db) }),
			db.NewClickHouseStore(func(db *sql.DB) { migrations.InitClickHouseMigrations(db) }),
			db.NewRedisStore(),
		),
	)
	router.Logger.Fatal(router.Start(":8080"))
}
