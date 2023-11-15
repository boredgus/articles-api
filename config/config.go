package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type config struct {
	MySQLUsername string `env:"MYSQL_USERNAME"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLDatabase string `env:"MYSQL_DATABASE"`
	DBContainer   string `env:"DB_CONTAINER"`
}

func LoadEnvFile(envFilePath string) {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Info("failed to load env file", err)
	}
}

func GetConfig() config {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		log.Error("failed to load env file", err)
	}
	return cfg
}

func InitConfig() {
	LoadEnvFile(".env")
	log.Infof("env config: %+v\n", GetConfig())
}
