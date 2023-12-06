package config

import (
	"errors"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type config struct {
	MySQLUsername        string `env:"MYSQL_USERNAME"`
	MySQLPassword        string `env:"MYSQL_PASSWORD"`
	MySQLDatabase        string `env:"MYSQL_DATABASE"`
	MaxOpenDBConnections int    `env:"MAX_OPEN_DB_CONNECTIONS"`
	MaxIdleDBConnections int    `env:"MAX_IDLE_DB_CONNECTIONS"`
	DBContainer          string `env:"DB_CONTAINER"`
}

func LoadEnvFile(envFilePath string) {
	_, err := os.Stat(envFilePath)
	if errors.Is(err, os.ErrNotExist) {
		logrus.Infof("%v file is not provided, skipping loading", envFilePath)
		return
	}
	if err = godotenv.Load(envFilePath); err != nil {
		logrus.Fatalf("failed to load %v file: %v", envFilePath, err)
	}
}

var cfg config
var configOnce sync.Once

func GetConfig() config {
	configOnce.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			logrus.Error("failed to load env file", err)
			cfg = config{}
		}
	})
	return cfg
}

func InitConfig() {
	LoadEnvFile(".env")
}
