// Articles API
//
// Spec Documentation for article service.
//
// Version: 1.0.0
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
	infrastructure "user-management/cmd/user-management/internal"
	"user-management/config"
)

func init() {
	config.InitConfig()
}

//go:generate swagger generate spec --scan-models --input=../../docs/init.json --output=../../docs/swagger.json
func main() {
	router := infrastructure.GetRouter(
		infrastructure.NewAppController(
			infrastructure.NewMySQLStore(),
		),
	)
	router.Logger.Fatal(router.Start(":8080"))
}
