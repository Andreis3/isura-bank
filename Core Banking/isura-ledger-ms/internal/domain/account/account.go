package account

import (
	"errors"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var (
	ErrInvalidAccountingType = errors.New("invalid accounting type")
	ErrEmptyExternalID       = errors.New("external id cannot be empty")
)

type AccountID string
type AccountingType string

const (
	Asset     AccountingType = "ASSET"
	Liability AccountingType = "LIABILITY"
	Revenue   AccountingType = "REVENUE"
	Expense   AccountingType = "EXPENSE"
)

func (a AccountingType) IsValid() bool {
	switch a {
	case Asset, Liability, Revenue, Expense:
		return true
	}
	return false
}

type Account struct {
	ID             AccountID
	ExternalID     string
	AccountingType AccountingType
	Balance        money.Money
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewAccount(accountID AccountID, externalID string, accountingType AccountingType, currency money.Currency) (*Account, error) {
	if !accountingType.IsValid() {
		return nil, ErrInvalidAccountingType
	}

	if externalID == "" {
		return nil, ErrEmptyExternalID
	}

	balance, err := money.NewMoney(0, currency)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Account{
		ID:             accountID,
		ExternalID:     externalID,
		AccountingType: accountingType,
		Balance:        balance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
