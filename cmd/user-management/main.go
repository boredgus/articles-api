package main

import (
	infrastructure "user-management/cmd/user-management/internal"
	"user-management/config"
)

func init() {
	config.InitConfig()
}

func main() {
	router := infrastructure.GetRouter(
		infrastructure.NewAppController(
			infrastructure.NewMySQLStore(),
		),
	)
	router.Logger.Fatal(router.Start(":8080"))
}
