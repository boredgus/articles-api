package main

import (
	"time"
	"user-management/config"
	"user-management/internal/infrastructure"

	"github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
}

func main() {
	router := infrastructure.GetRouter()
	db, err := infrastructure.NewMySQLStore()
	if err != nil {
		logrus.Warnf("failed to init db: %v", err)
		return
	}

	time.Sleep(10 * time.Second)
	rows, err := db.Query("select * from user;")
	if err != nil {
		logrus.Infof("failed to fetch users: %v", err)
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		logrus.Infof("failed to get columns: %v", err)
		return
	}

	logrus.Infof("columns: %v", cols)

	router.Logger.Fatal(router.Start(":8080"))
}
