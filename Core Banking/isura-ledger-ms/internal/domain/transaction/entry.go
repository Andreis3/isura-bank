package transaction

import (
	"errors"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var (
	ErrInvalidDirection    = errors.New("invalid direction")
	ErrNegativeAmountValue = errors.New("amount cannot be negative")
	ErrAmountEqualZero     = errors.New("amount cannot be zero")
)

type Direction string

const (
	Credit Direction = "CREDIT"
	Debit  Direction = "DEBIT"
)

func (d Direction) IsValid() bool {
	switch d {
	case Credit, Debit:
		return true
	}
	return false
}

type Entry struct {
	ID             EntryID
	IdempotencyKey string
	Direction      Direction
	Amount         money.Money
	AccountID      AccountID
	TransactionID  TransactionID
	CreatedAt      time.Time
}

func NewEntry(entryID EntryID,
	idempotencyKey string,
	direction Direction,
	amount money.Money,
	accountID AccountID,
	transactionID TransactionID) (*Entry, error) {
	if !direction.IsValid() {
		return nil, ErrInvalidDirection
	}

	if amount.IsZero() {
		return nil, ErrAmountEqualZero
	}

	if amount.IsNegative() {
		return nil, ErrNegativeAmountValue
	}

	return &Entry{
		ID:             entryID,
		IdempotencyKey: idempotencyKey,
		Direction:      direction,
		Amount:         amount,
		AccountID:      accountID,
		TransactionID:  transactionID,
		CreatedAt:      time.Now(),
	}, nil
}
