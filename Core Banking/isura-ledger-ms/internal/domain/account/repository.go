package account

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Repository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id AccountID) (*Account, error)
	FindByExternalID(ctx context.Context, externalID string) (*Account, error)
	UpdateBalance(ctx context.Context, accountID AccountID, balance money.Money) error
	FindBalanceByID(ctx context.Context, accountID AccountID) (money.Money, error)
	FindBalanceForUpdateByID(ctx context.Context, accountID AccountID) (money.Money, error)
}
