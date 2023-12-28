package gateways

import (
	"database/sql"
)

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type Store interface {
	Query(query string, args ...any) (Rows, error)
	Stats() sql.DBStats
}

type CacheStore interface {
	Get(key string, value any) error
	Set(key string, value any) error
}
