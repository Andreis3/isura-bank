package model

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction struct {
	ID             pgtype.Text
	IdempotencyKey pgtype.Text
	Status         pgtype.Text
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

func ToTransactionModel(domain *transaction.Transaction) Transaction {
	return Transaction{
		ID: pgtype.Text{
			String: string(domain.ID),
			Valid:  true,
		},
		IdempotencyKey: pgtype.Text{
			String: domain.IdempotencyKey,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: string(domain.Status),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  domain.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  domain.UpdatedAt,
			Valid: true,
		},
	}
}

func ToTransactionDomain(model Transaction, entries []*transaction.Entry) *transaction.Transaction {
	return &transaction.Transaction{
		ID:             transaction.TransactionID(model.ID.String),
		IdempotencyKey: model.IdempotencyKey.String,
		Status:         transaction.TransactionStatus(model.Status.String),
		Entries:        entries,
		CreatedAt:      model.CreatedAt.Time,
		UpdatedAt:      model.UpdatedAt.Time,
	}
}
