package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
)

type ObservabilityOutboxRepo struct {
	repo   outbox.Repository
	metric application.Metrics
	tracer application.Tracer
}

func NewObservabilityOutboxRepo(repo outbox.Repository, metric application.Metrics, tracer application.Tracer) *ObservabilityOutboxRepo {
	return &ObservabilityOutboxRepo{
		repo:   repo,
		metric: metric,
		tracer: tracer,
	}
}

func (r *ObservabilityOutboxRepo) Save(ctx context.Context, outbox *outbox.Outbox) error {
	ctx, span := r.tracer.Start(ctx, "OutboxRepository.Save")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"outbox",
			"save",
			float64(time.Since(start).Milliseconds()))
	}()
	err := r.repo.Save(ctx, outbox)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityOutboxRepo) FindAllByStatusForUpdateSkipLocked(ctx context.Context, status outbox.StatusOutbox, limit int) ([]*outbox.Outbox, error) {
	ctx, span := r.tracer.Start(ctx, "OutboxRepository.FindAllByStatusForUpdateSkipLocked")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"outbox",
			"find_all_Status_for_skip_locked",
			float64(time.Since(start).Milliseconds()))
	}()

	outboxes, err := r.repo.FindAllByStatusForUpdateSkipLocked(ctx, status, limit)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return outboxes, nil
}

func (r *ObservabilityOutboxRepo) UpdateOutboxData(ctx context.Context, outboxID outbox.OutboxID, data outbox.UpdateOutboxData) error {
	ctx, span := r.tracer.Start(ctx, "OutboxRepository.UpdateOutboxData")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"outbox",
			"update_outbox_data",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.UpdateOutboxData(ctx, outboxID, data)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
