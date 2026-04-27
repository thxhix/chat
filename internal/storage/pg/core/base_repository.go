package core

import (
	"context"
	"database/sql"
	"github.com/thxhix/chat/internal/storage/pg/core/tx_manager"
)

type BaseRepository struct {
	DB        *sql.DB
	txManager *tx_manager.TXManager
}

func NewBaseRepository(db *sql.DB, txManager *tx_manager.TXManager) BaseRepository {
	return BaseRepository{
		DB:        db,
		txManager: txManager,
	}
}

func (r *BaseRepository) GetExecutor(ctx context.Context) tx_manager.DBExecutor {
	return r.txManager.GetExecutor(ctx, r.DB)
}
