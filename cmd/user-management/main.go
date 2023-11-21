package main

import (
	infrastructure "user-management/cmd/user-management/internal"
	"user-management/config"

	"github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
}

func main() {
	db, err := infrastructure.NewMySQLStore()
	if err != nil {
		logrus.Warnf("failed to init db: %v", err)
		return
	}

	router := infrastructure.GetRouter(infrastructure.NewAppController(db))
	router.Logger.Fatal(router.Start(":8080"))
}
