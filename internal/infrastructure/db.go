package infrastructure

import (
	"database/sql"
	"fmt"
	"time"
	"user-management/config"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Query()
}

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore() (MySQLStore, error) {
	// uri, err := NewMySQLURI().Generate()
	// if err != nil {
	// return MySQLStore{}, err
	// }

	config := config.GetConfig()
	defaultConfig := mysql.NewConfig()
	cnfg := mysql.Config{
		User:                 config.MySQLUsername,
		Passwd:               config.MySQLPassword,
		Net:                  "tcp",
		Addr:                 config.DBContainer,
		DBName:               config.MySQLDatabase,
		Collation:            defaultConfig.Collation,
		Loc:                  defaultConfig.Loc,
		MaxAllowedPacket:     defaultConfig.MaxAllowedPacket,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}
	logrus.Infof("URI = %v", cnfg.FormatDSN())
	db, err := sql.Open("mysql", cnfg.FormatDSN())

	if err != nil {
		return MySQLStore{}, fmt.Errorf("failed to connect db: %v", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	time.Sleep(10 * time.Second)
	err = db.Ping()
	logrus.Infof("after ping: %v", err)

	return MySQLStore{db: db}, nil
}

func (s MySQLStore) Query(q string, args ...any) (*sql.Rows, error) {
	return s.db.Query(q, args...)
}
