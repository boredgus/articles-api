package gateways

import "database/sql"

type Store interface {
	Query(query string, args ...any) (*sql.Rows, error)
}
