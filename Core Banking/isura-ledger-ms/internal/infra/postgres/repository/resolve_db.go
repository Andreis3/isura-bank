package repository

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
)

func resolveDB(ctx context.Context, database database.Querier) database.Querier {
	if tx, ok := database.ExtractTx(ctx); ok {
		return tx
	}
	return database
}
