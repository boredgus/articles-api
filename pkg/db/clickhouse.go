package db

import (
	"a-article/config"
	"a-article/internal/gateways"
	"database/sql"
	"fmt"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sirupsen/logrus"
)

var clickhouseDB *sql.DB
var clickhouseOnce sync.Once

func NewClickHouseStore(handlers ...func(db *sql.DB)) gateways.Store {
	clickhouseOnce.Do(func() {
		cfg := config.GetConfig()
		db := clickhouse.OpenDB(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%v:%v", cfg.ClickhouseContainer, cfg.ClickhousePort)},
			Auth: clickhouse.Auth{
				Database: cfg.ClickhouseDatabase,
				Username: cfg.ClickhouseUsername,
				Password: cfg.ClickhousePassword,
			},
			Debug:    cfg.Mode == config.Dev,
			Protocol: clickhouse.HTTP,
		})
		db.SetMaxOpenConns(cfg.MaxOpenDBConnections)
		db.SetMaxIdleConns(cfg.MaxIdleDBConnections)
		for _, handler := range handlers {
			handler(db)
		}
		clickhouseDB = db
	})
	return ClickhouseStore{db: clickhouseDB}
}

type ClickhouseStore struct {
	db *sql.DB
}

func (s ClickhouseStore) Query(query string, args ...any) (gateways.Rows, error) {
	logrus.Infof("> clickhouse: %+v", s.Stats())
	return s.db.Query(query, args...)
}

func (s ClickhouseStore) Stats() sql.DBStats {
	return s.db.Stats()
}
