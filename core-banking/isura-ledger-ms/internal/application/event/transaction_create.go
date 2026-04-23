package event

import "time"

type TransactionCreated struct {
	TransactionID   string    `json:"transaction_id"`
	IdempotencyKey  string    `json:"idempotency_key"`
	DebitAccountID  string    `json:"debit_account_id"`
	CreditAccountID string    `json:"credit_account_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	OccurredAt      time.Time `json:"occurred_at"`
}
