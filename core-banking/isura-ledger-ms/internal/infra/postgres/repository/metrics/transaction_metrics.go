package metrics

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

type MetricsTransactionRepo struct {
	repo   transaction.Repository
	metric application.Metrics
}

func NewMetricsTransactionRepo(repo transaction.Repository, metric application.Metrics) *MetricsTransactionRepo {
	return &MetricsTransactionRepo{
		repo:   repo,
		metric: metric,
	}
}

func (r *MetricsTransactionRepo) Save(ctx context.Context, data *transaction.Transaction) error {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"save",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.Save(ctx, data)

	if err != nil {
		return err
	}

	return nil
}

func (r *MetricsTransactionRepo) FindByID(ctx context.Context, transactionID transaction.TransactionID) (*transaction.Transaction, error) {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"find_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	transaction, err := r.repo.FindByID(ctx, transactionID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *MetricsTransactionRepo) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*transaction.Transaction, error) {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"find_by_idempotency_key",
			float64(time.Since(start).Milliseconds()))
	}()

	transaction, err := r.repo.FindByIdempotencyKey(ctx, idempotencyKey)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *MetricsTransactionRepo) ExistsByIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error) {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"exists_by_idempotency_key",
			float64(time.Since(start).Milliseconds()))
	}()

	exists, err := r.repo.ExistsByIdempotencyKey(ctx, idempotencyKey)

	if err != nil {
		return false, err
	}

	return exists, nil
}
