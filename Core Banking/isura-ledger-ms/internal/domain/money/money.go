package money

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCurrency  = errors.New("invalid currency")
	ErrCurrencyMismatch = errors.New("currencies must be the same")
	ErrNegativeAmount   = errors.New("amount cannot be negative")
)

type Currency string

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

func (c Currency) IsValid() bool {
	switch c {
	case BRL, USD, EUR:
		return true
	}
	return false
}

type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}

	if !currency.IsValid() {
		return Money{}, ErrInvalidCurrency
	}
	return Money{amount, currency}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return Money{m.amount + other.amount, m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{m.amount - other.amount, m.currency}, nil
}
func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) Equal(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) IsSufficientBalance(other Money) bool {
	return m.amount >= other.amount && m.currency == other.currency
}

func (m Money) String() string {
	units := m.amount / 100
	cents := m.amount % 100
	return fmt.Sprintf("%d.%02d %s", units, cents, m.currency)
}
