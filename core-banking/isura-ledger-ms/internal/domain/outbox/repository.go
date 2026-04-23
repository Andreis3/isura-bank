package outbox

import (
	"context"
	"time"
)

type UpdateOutboxData struct {
	Status        StatusOutbox
	Attempts      int
	LastAttemptAt *time.Time
	PublishedAt   *time.Time
}
type Repository interface {
	Save(ctx context.Context, outbox *Outbox) error
	FindAllByStatusForUpdateSkipLocked(ctx context.Context, status StatusOutbox, limit int) ([]*Outbox, error)
	UpdateOutboxData(ctx context.Context, outboxID OutboxID, data UpdateOutboxData) error
}
