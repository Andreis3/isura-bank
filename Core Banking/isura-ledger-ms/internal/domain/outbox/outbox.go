package outbox

import (
	"errors"
	"time"
)

const MaxAttempts = 3

var (
	ErrMaxAttemptsExceeded = errors.New("max attempts exceeded")
	ErrTransitionStatus    = errors.New("transition status")
)

type OutboxID string

type AggregateType string

type EventType string

type StatusOutbox string

const (
	Transaction AggregateType = "TRANSACTION"
)

const (
	TransactionCreated EventType = "ledger.transaction.created"
	BalanceUpdated     EventType = "ledger.balance.updated"
	EntryCreated       EventType = "ledger.entry.created"
)

const (
	Pending StatusOutbox = "PENDING"
	Failed  StatusOutbox = "FAILED"
	Success StatusOutbox = "SUCCESS"
)

type StateMachineStatus map[StatusOutbox][]StatusOutbox

var validStateMachine = StateMachineStatus{
	Pending: []StatusOutbox{Failed, Success},
	Failed:  []StatusOutbox{Pending},
	Success: []StatusOutbox{},
}

func (s StatusOutbox) IsValid() bool {
	switch s {
	case Pending, Failed, Success:
		return true
	}
	return false
}

func (s StateMachineStatus) CanTransition(from, to StatusOutbox) bool {
	allowed, ok := s[from]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}

type Outbox struct {
	ID            OutboxID
	AggregateID   string
	AggregateType AggregateType
	EventType     EventType
	Payload       []byte
	Status        StatusOutbox
	Attempts      int
	LastAttemptAt *time.Time
	CreatedAt     time.Time
	PublishedAt   *time.Time
}

func NewOutbox(outboxID OutboxID,
	aggregateID string,
	aggregateType AggregateType,
	eventType EventType,
	payload []byte) *Outbox {
	return &Outbox{
		ID:            outboxID,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		Payload:       payload,
		Status:        Pending,
		Attempts:      0,
		CreatedAt:     time.Now(),
	}
}

func (o *Outbox) Publish() error {
	if !validStateMachine.CanTransition(o.Status, Success) {
		return ErrTransitionStatus
	}
	o.Status = Success
	o.PublishedAt = new(time.Now())
	return nil
}

func (o *Outbox) MarkFailed() error {
	if !validStateMachine.CanTransition(o.Status, Failed) {
		return ErrTransitionStatus
	}

	o.Status = Failed
	o.Attempts++
	o.LastAttemptAt = new(time.Now())

	return nil
}

func (o *Outbox) Retry() error {
	if o.Attempts >= MaxAttempts {
		return ErrMaxAttemptsExceeded
	}

	if !validStateMachine.CanTransition(o.Status, Pending) {
		return ErrTransitionStatus
	}

	o.Status = Pending

	return nil
}
