package transaction

import (
	"errors"
	"time"
)

var (
	ErrInvalidMaxEntries        = errors.New("maximum entries exceeded")
	ErrDuplicateEntryDirection  = errors.New("duplicate entry direction")
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
	ErrInvalidDifferentAmount   = errors.New("different amount")
	ErrTransactionNotFound      = errors.New("transaction not found")
)

type TransactionStatus string

const (
	Pending   TransactionStatus = "PENDING"
	Completed TransactionStatus = "COMPLETED"
	Failed    TransactionStatus = "FAILED"
)

func (t TransactionStatus) IsValid() bool {
	switch t {
	case Pending, Completed, Failed:
		return true
	}
	return false
}

type Transaction struct {
	ID             TransactionID
	IdempotencyKey string
	Status         TransactionStatus
	Entries        []*Entry
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewTransaction(transactionID TransactionID, idempotencyKey string) *Transaction {
	createDate := time.Now()
	return &Transaction{
		ID:             transactionID,
		IdempotencyKey: idempotencyKey,
		Entries:        make([]*Entry, 0, 2),
		Status:         Pending,
		CreatedAt:      createDate,
		UpdatedAt:      createDate,
	}
}

func (t *Transaction) AddEntry(entry *Entry) error {
	if len(t.Entries) >= 2 {
		return ErrInvalidMaxEntries
	}

	if len(t.Entries) == 1 {

		if t.Entries[0].Direction == entry.Direction {
			return ErrDuplicateEntryDirection
		}

		if !t.Entries[0].Amount.Equal(entry.Amount) {
			return ErrInvalidDifferentAmount
		}
	}

	t.Entries = append(t.Entries, entry)

	return nil
}

func (t *Transaction) Complete() error {
	if !validStateMachine.CanTransition(t.Status, Completed) {
		return ErrInvalidTransactionStatus
	}

	t.Status = Completed
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) Fail() error {
	if !validStateMachine.CanTransition(t.Status, Failed) {
		return ErrInvalidTransactionStatus
	}

	t.Status = Failed
	t.UpdatedAt = time.Now()
	return nil
}
