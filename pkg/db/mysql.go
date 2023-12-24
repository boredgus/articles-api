package db

import (
	"database/sql"
	"sync"
	"time"
	"user-management/config"
	"user-management/internal/gateways"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func getConfig() *mysql.Config {
	env := config.GetConfig()
	config := mysql.NewConfig()

	config.User = env.MySQLUsername
	config.Passwd = env.MySQLPassword
	config.Net = "tcp"
	config.Addr = env.MySQLContainer
	config.DBName = env.MySQLDatabase
	config.ParseTime = true
	config.MultiStatements = true
	config.InterpolateParams = true

	return config
}

var mysqlOnce sync.Once
var mysqlDB *sql.DB

const DBConnectionAttempts = 5

func NewMySQLStore(handlers ...func(db *sql.DB)) gateways.Store {
	mysqlOnce.Do(func() {
		var db *sql.DB
		var err error
		for i := 0; i < DBConnectionAttempts; i++ {
			db, err = sql.Open("mysql", getConfig().FormatDSN())
			if err != nil {
				logrus.Warnf("db connection (attempt #%v) failed: %v", i+1, err)
				time.Sleep(10 * time.Second)
				continue
			}
			break
		}
		if err != nil {
			logrus.Fatalf("failed to connect to DB: %v", err)
		}

		cfg := config.GetConfig()
		db.SetMaxOpenConns(cfg.MaxOpenDBConnections)
		db.SetMaxIdleConns(cfg.MaxIdleDBConnections)
		mysqlDB = db
		for _, handler := range handlers {
			handler(mysqlDB)
		}
	})

	return MySQLStore{db: mysqlDB}
}

type MySQLStore struct {
	db *sql.DB
}

func (s MySQLStore) Query(q string, args ...any) (gateways.Rows, error) {
	logrus.Infof("> mysql: %+v\n", s.Stats())
	return s.db.Query(q, args...)
}
func (s MySQLStore) Stats() sql.DBStats {
	return s.db.Stats()
}
