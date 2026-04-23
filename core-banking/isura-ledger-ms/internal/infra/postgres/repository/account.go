package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/model"
)

type AccountRepository struct {
	db database.Querier
}

func NewAccountRepository(db database.Querier) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Save(ctx context.Context, account *account.Account) error {
	db := resolveDB(ctx, r.db)

	query := `INSERT INTO 
    accounts (id, external_id, account_type, balance, currency, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	accountModel := model.ToAccountModel(account)

	_, err := db.Exec(ctx, query,
		accountModel.ID,
		accountModel.ExternalID,
		accountModel.AccountType,
		accountModel.Balance,
		accountModel.Currency,
		accountModel.CreatedAt,
		accountModel.UpdatedAt,
	)

	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT 
    	id, external_id, account_type, balance, currency, created_at, updated_at 
		FROM accounts WHERE id = $1`

	var accountModel model.Account
	err := db.QueryRow(ctx, query, id).Scan(
		&accountModel.ID,
		&accountModel.ExternalID,
		&accountModel.AccountType,
		&accountModel.Balance,
		&accountModel.Currency,
		&accountModel.CreatedAt,
		&accountModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}

	account, err := model.ToAccountDomain(accountModel)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) FindByExternalID(ctx context.Context, externalID string) (*account.Account, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT 
    	id, external_id, account_type, balance, currency, created_at, updated_at 
		FROM accounts WHERE external_id = $1`

	var accountModel model.Account
	err := db.QueryRow(ctx, query, externalID).Scan(
		&accountModel.ID,
		&accountModel.ExternalID,
		&accountModel.AccountType,
		&accountModel.Balance,
		&accountModel.Currency,
		&accountModel.CreatedAt,
		&accountModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}

	account, err := model.ToAccountDomain(accountModel)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID account.AccountID, balance money.Money) error {
	db := resolveDB(ctx, r.db)
	query := `UPDATE accounts SET balance = $1, currency = $2, updated_at = $3 WHERE id = $4`

	_, err := db.Exec(ctx, query,
		pgtype.Int8{Int64: balance.Amount(), Valid: true},
		pgtype.Text{String: string(balance.Currency()), Valid: true},
		pgtype.Timestamptz{Time: time.Now(), Valid: true},
		pgtype.Text{String: string(accountID), Valid: true},
	)

	return err
}

func (r *AccountRepository) FindBalanceByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT balance, currency FROM accounts WHERE id = $1`
	var accountModel model.Account

	err := db.QueryRow(ctx, query, accountID).Scan(
		&accountModel.Balance,
		&accountModel.Currency,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return money.Money{}, account.ErrAccountNotFound
		}
		return money.Money{}, err
	}

	newMoney, err := money.NewMoney(accountModel.Balance.Int64, money.Currency(accountModel.Currency.String))
	if err != nil {
		return money.Money{}, err
	}

	return newMoney, nil
}

func (r *AccountRepository) FindBalanceForUpdateByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT balance, currency FROM accounts WHERE id = $1 FOR UPDATE`
	var accountModel model.Account

	err := db.QueryRow(ctx, query, accountID).Scan(
		&accountModel.Balance,
		&accountModel.Currency,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return money.Money{}, account.ErrAccountNotFound
		}
		return money.Money{}, err
	}

	newMoney, err := money.NewMoney(accountModel.Balance.Int64, money.Currency(accountModel.Currency.String))
	if err != nil {
		return money.Money{}, err
	}

	return newMoney, nil
}
