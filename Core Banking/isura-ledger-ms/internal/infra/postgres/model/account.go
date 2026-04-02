package model

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

type Account struct {
	ID          pgtype.Text
	ExternalID  pgtype.Text
	AccountType pgtype.Text
	Balance     pgtype.Int8
	Currency    pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func ToAccountModel(domain *account.Account) Account {
	return Account{
		ID: pgtype.Text{
			String: string(domain.ID),
			Valid:  true,
		},
		ExternalID: pgtype.Text{
			String: domain.ExternalID,
			Valid:  true,
		},
		AccountType: pgtype.Text{
			String: string(domain.AccountType),
			Valid:  true,
		},
		Balance: pgtype.Int8{
			Int64: domain.Balance.Amount(),
			Valid: true,
		},
		Currency: pgtype.Text{
			String: string(domain.Balance.Currency()),
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

func ToAccountDomain(model Account) (*account.Account, error) {
	balance, err := money.NewMoney(model.Balance.Int64, money.Currency(model.Currency.String))
	if err != nil {
		return nil, err
	}

	return &account.Account{
		ID:          account.AccountID(model.ID.String),
		ExternalID:  model.ExternalID.String,
		Balance:     balance,
		AccountType: account.AccountType(model.AccountType.String),
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}, nil
}
