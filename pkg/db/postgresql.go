package db

import (
	"database/sql"
	"fmt"
	"sync"
	"user-management/config"
	"user-management/internal/gateways"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var postgresOnce sync.Once
var postgresDB *sql.DB

func NewPostrgreSQLStore(handlers ...func(db *sql.DB)) gateways.Store {
	postgresOnce.Do(func() {
		env := config.GetConfig()
		db, err := sql.Open("postgres",
			fmt.Sprintf("user=%s password=%s host=%s dbname=%s search_path=%s sslmode=disable",
				env.PostgreSQLUser,
				env.PostgreSQLPassword,
				env.PostgreSQLContainer,
				env.PostgreSQLDatabase,
				env.PostgreSQLDatabase,
			))
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres db: %w", err))
		}

		postgresDB = db
	})
	for _, handler := range handlers {
		handler(postgresDB)
	}
	return &PostrgreSQLStore{db: postgresDB}
}

type PostrgreSQLStore struct {
	db *sql.DB
}

func (s PostrgreSQLStore) Query(query string, args ...any) (gateways.Rows, error) {
	logrus.Infof("> postgresql: %+v\n", s.Stats())
	return s.db.Query(query, args...)
}
func (s PostrgreSQLStore) Stats() sql.DBStats {
	return s.db.Stats()
}
