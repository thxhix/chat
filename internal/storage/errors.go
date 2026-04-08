package storage

import "errors"

var (
	ErrNoPostgresConnection = errors.New("no postgres connection")
)
