// Articles API
//
// Spec Documentation for article service.
//
// Version: 1.2.1
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
	"a-article/cmd/articles-api/internal"
	"a-article/config"
	"a-article/pkg/db"
	"a-article/pkg/msgbroker"
	migrations "a-article/sql"
	"database/sql"
)

func init() {
	config.InitConfig()
}

//go:generate swagger generate spec --scan-models --input=../../docs/init.json --output=../../docs/swagger.json
func main() {
	router := internal.GetRouter(
		internal.NewAppController(
			internal.AppParams{
				MainStore:     db.NewPostrgreSQLStore(func(db *sql.DB) { migrations.InitPostgreSQLMigrations(db) }),
				StatsStore:    db.NewClickHouseStore(func(db *sql.DB) { migrations.InitClickHouseMigrations(db) }),
				CacheStore:    db.NewRedisStore(),
				MessageBroker: msgbroker.NewRabbitMQ(),
			},
		),
	)
	router.Logger.Fatal(router.Start(":8080"))
}
