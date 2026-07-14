package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxManager interface {
	WithTx(ctx context.Context, txFn func(tx *sqlx.Tx) error) error
	GetDB() *sqlx.DB
}

type txManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) TxManager {
	return &txManager{
		db: db,
	}
}

func (s *txManager) WithTx(ctx context.Context, txFn func(tx *sqlx.Tx) error) error {
	op := "TxManager.WithTx"
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	if err := txFn(tx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

func (s *txManager) GetDB() *sqlx.DB {
	return s.db
}
