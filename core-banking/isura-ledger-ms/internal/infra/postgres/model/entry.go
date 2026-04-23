package model

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/jackc/pgx/v5/pgtype"
)

type Entry struct {
	ID             pgtype.Text
	IdempotencyKey pgtype.Text
	Direction      pgtype.Text
	Amount         pgtype.Int8
	Currency       pgtype.Text
	AccountID      pgtype.Text
	TransactionID  pgtype.Text
	CreatedAt      pgtype.Timestamptz
}

func ToEntryModel(domain *transaction.Entry) Entry {
	return Entry{
		ID: pgtype.Text{
			String: string(domain.ID),
			Valid:  true,
		},
		IdempotencyKey: pgtype.Text{
			String: domain.IdempotencyKey,
			Valid:  true,
		},
		Direction: pgtype.Text{
			String: string(domain.Direction),
			Valid:  true,
		},
		Amount: pgtype.Int8{
			Int64: domain.Amount.Amount(),
			Valid: true,
		},
		Currency: pgtype.Text{
			String: string(domain.Amount.Currency()),
			Valid:  true,
		},
		AccountID: pgtype.Text{
			String: string(domain.AccountID),
			Valid:  true,
		},
		TransactionID: pgtype.Text{
			String: string(domain.TransactionID),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  domain.CreatedAt,
			Valid: true,
		},
	}
}

func ToEntryDomain(model Entry) (*transaction.Entry, error) {
	amount, err := money.NewMoney(model.Amount.Int64, money.Currency(model.Currency.String))
	if err != nil {
		return nil, err
	}
	return &transaction.Entry{
		ID:             transaction.EntryID(model.ID.String),
		IdempotencyKey: model.IdempotencyKey.String,
		TransactionID:  transaction.TransactionID(model.TransactionID.String),
		AccountID:      transaction.AccountID(model.AccountID.String),
		Amount:         amount,
		Direction:      transaction.Direction(model.Direction.String),
		CreatedAt:      model.CreatedAt.Time,
	}, nil
}
