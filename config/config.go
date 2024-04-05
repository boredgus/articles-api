package config

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Mode int

const (
	Dev Mode = iota
	Production
)

type config struct {
	Mode         Mode   `env:"MODE"`
	JWTSecretKey string `env:"JWT_SECRET_KEY"`

	BrevoAPIKey  string `env:"BREVO_API_KEY"`
	SMTPUsername string `env:"SMTP_USERNAME"`

	MySQLContainer       string `env:"MYSQL_CONTAINER"`
	MySQLUsername        string `env:"MYSQL_USERNAME"`
	MySQLPassword        string `env:"MYSQL_PASSWORD"`
	MySQLDatabase        string `env:"MYSQL_DATABASE"`
	MaxOpenDBConnections int    `env:"MAX_OPEN_DB_CONNECTIONS"`
	MaxIdleDBConnections int    `env:"MAX_IDLE_DB_CONNECTIONS"`

	ClickhousePort      int    `env:"CLICKHOUSE_PORT"`
	ClickhouseContainer string `env:"CLICKHOUSE_CONTAINER"`
	ClickhouseUsername  string `env:"CLICKHOUSE_USERNAME"`
	ClickhousePassword  string `env:"CLICKHOUSE_PASSWORD"`
	ClickhouseDatabase  string `env:"CLICKHOUSE_DATABASE"`

	RedisPort              string `env:"REDIS_PORT"`
	RedisContainer         string `env:"REDIS_CONTAINER"`
	RedisPassword          string `env:"REDIS_PASSWORD"`
	RedisExpirationMinutes int64  `env:"REDIS_EXPIRATION_MINUTES"`

	PostgreSQLContainer string `env:"POSTGRESQL_CONTAINER"`
	PostgreSQLUser      string `env:"POSTGRESQL_USER"`
	PostgreSQLPassword  string `env:"POSTGRESQL_PASSWORD"`
	PostgreSQLDatabase  string `env:"POSTGRESQL_DB"`

	RabbitMQContainer string `env:"RABBITMQ_CONTAINER"`
	RabbitMQPort      string `env:"RABBITMQ_PORT"`
	RabbitMQUser      string `env:"RABBITMQ_DEFAULT_USER"`
	RabbitMQPass      string `env:"RABBITMQ_DEFAULT_PASS"`
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

func InitConfig() {
	LoadEnvFile(".env")
	if err := env.Parse(&cfg); err != nil {
		logrus.Error("failed to load env file", err)
		cfg = config{}
	}
}

func GetConfig() config {
	return cfg
}
