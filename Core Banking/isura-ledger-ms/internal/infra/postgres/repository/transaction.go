package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/model"
)

type TransactionRepository struct {
	db database.Querier
}

func NewTransaction(db database.Querier) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, data *transaction.Transaction) error {
	batch := pgx.Batch{}

	transactionModel := model.ToTransactionModel(data)

	batch.Queue(`
		INSERT INTO transactions 
			(id, idempotency_key, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		transactionModel.ID,
		transactionModel.IdempotencyKey,
		transactionModel.Status,
		transactionModel.CreatedAt,
		transactionModel.UpdatedAt)

	for _, entry := range data.Entries {
		entryModel := model.ToEntryModel(entry)
		batch.Queue(`
			INSERT INTO entries 
				(id, idempotency_key, direction, amount, currency, account_id, transaction_id, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			entryModel.ID,
			entryModel.IdempotencyKey,
			entryModel.Direction,
			entryModel.Amount,
			entryModel.Currency,
			entryModel.AccountID,
			entryModel.TransactionID,
			entryModel.CreatedAt)
	}

	db := resolveDB(ctx, r.db)

	results := db.SendBatch(ctx, &batch)
	defer results.Close()

	if _, err := results.Exec(); err != nil {
		return err
	}

	for range data.Entries {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, transactionID transaction.TransactionID) (*transaction.Transaction, error) {
	db := resolveDB(ctx, r.db)

	batch := pgx.Batch{}

	queryTransaction := `
	SELECT id, idempotency_key, status, created_at, updated_at
	FROM transactions
	WHERE id = $1`

	queryEntries := `
	SELECT id, idempotency_key, direction, amount, currency, account_id, transaction_id, created_at
	FROM entries
	WHERE transaction_id = $1`

	batch.Queue(queryTransaction, transactionID)
	batch.Queue(queryEntries, transactionID)

	results := db.SendBatch(ctx, &batch)
	defer results.Close()

	// first result transaction
	transactioRow := results.QueryRow()
	var transactionModel model.Transaction
	if err := transactioRow.Scan(
		&transactionModel.ID,
		&transactionModel.IdempotencyKey,
		&transactionModel.Status,
		&transactionModel.CreatedAt,
		&transactionModel.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction.ErrTransactionNotFound
		}
		return nil, err
	}

	// second result entries
	entryRows, err := results.Query()
	if err != nil {
		return nil, err
	}
	defer entryRows.Close()

	var entries []*transaction.Entry
	for entryRows.Next() {
		var entryModel model.Entry
		if err := entryRows.Scan(
			&entryModel.ID,
			&entryModel.IdempotencyKey,
			&entryModel.Direction,
			&entryModel.Amount,
			&entryModel.Currency,
			&entryModel.AccountID,
			&entryModel.TransactionID,
			&entryModel.CreatedAt,
		); err != nil {
			return nil, err
		}

		entry, err := model.ToEntryDomain(entryModel)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := entryRows.Err(); err != nil {
		return nil, err
	}

	return model.ToTransactionDomain(transactionModel, entries), nil

}

func (r *TransactionRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*transaction.Transaction, error) {
	db := resolveDB(ctx, r.db)

	query := `
		SELECT 
		    t.id, t.idempotency_key, t.status, t.created_at, t.updated_at,
		    e.id, e.idempotency_key, e.direction, e.amount, e.currency, 
		    e.account_id, e.transaction_id, e.created_at
		FROM transactions t
		JOIN entries e ON e.transaction_id = t.id
		WHERE t.idempotency_key = $1
	`

	rows, err := db.Query(ctx, query, idempotencyKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionModel model.Transaction
	var entryModels []model.Entry
	found := false

	for rows.Next() {
		var entry model.Entry
		if err := rows.Scan(
			&transactionModel.ID,
			&transactionModel.IdempotencyKey,
			&transactionModel.Status,
			&transactionModel.CreatedAt,
			&transactionModel.UpdatedAt,
			&entry.ID,
			&entry.IdempotencyKey,
			&entry.Direction,
			&entry.Amount,
			&entry.Currency,
			&entry.AccountID,
			&entry.TransactionID,
			&entry.CreatedAt,
		); err != nil {
			return nil, err
		}

		entryModels = append(entryModels, entry)
		found = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, transaction.ErrTransactionNotFound
	}

	var entries []*transaction.Entry
	for _, entryModel := range entryModels {
		entry, err := model.ToEntryDomain(entryModel)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return model.ToTransactionDomain(transactionModel, entries), nil

}

func (r *TransactionRepository) ExistsByIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error) {
	db := resolveDB(ctx, r.db)
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM transactions WHERE idempotency_key = $1)
	`, idempotencyKey).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
