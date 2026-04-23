package uow

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBeginTransaction    = errors.New("error opening transaction")
	ErrCommitTransaction   = errors.New("error committing transaction")
	ErrRollbackTransaction = errors.New("error rolling back transaction")
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBeginTransaction, err)
	}

	if err := fn(database.WithTx(ctx, tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("%w: %v | rollback error: %v", ErrRollbackTransaction, err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrCommitTransaction, err)
	}

	return nil
}
