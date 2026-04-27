package tx_manager

import (
	"context"
	"database/sql"
	"fmt"
)

type TXKey string

type ITXManager interface {
	GetExecutor(ctx context.Context, defaultDB *sql.DB) DBExecutor
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type TXManager struct {
	Key TXKey
	db  *sql.DB
}

func NewTXManager(key TXKey, db *sql.DB) *TXManager {
	return &TXManager{
		Key: key,
		db:  db,
	}
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (m *TXManager) GetExecutor(ctx context.Context, defaultDB *sql.DB) DBExecutor {
	if tx, ok := ctx.Value(m.Key).(*sql.Tx); ok {
		return tx
	}
	return defaultDB
}

func (m *TXManager) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	ctxWithTx := context.WithValue(ctx, m.Key, tx)

	err = tFunc(ctxWithTx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction failed: %v (rollback error: %v)", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
