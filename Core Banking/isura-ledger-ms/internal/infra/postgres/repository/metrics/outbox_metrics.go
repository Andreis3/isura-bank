package metrics

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
)

type MetricsOutboxRepo struct {
	repo   outbox.Repository
	metric application.Metrics
}

func NewMetricsOutboxRepo(repo outbox.Repository, metric application.Metrics) *MetricsOutboxRepo {
	return &MetricsOutboxRepo{
		repo:   repo,
		metric: metric,
	}
}

func (r *MetricsOutboxRepo) Save(ctx context.Context, outbox *outbox.Outbox) error {
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
		return err
	}

	return nil
}

func (r *MetricsOutboxRepo) FindAllByStatusForUpdateSkipLocked(ctx context.Context, status outbox.StatusOutbox, limit int) ([]*outbox.Outbox, error) {
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
		return nil, err
	}

	return outboxes, nil
}

func (r *MetricsOutboxRepo) UpdateOutboxData(ctx context.Context, outboxID outbox.OutboxID, data outbox.UpdateOutboxData) error {
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
		return err
	}

	return nil
}
