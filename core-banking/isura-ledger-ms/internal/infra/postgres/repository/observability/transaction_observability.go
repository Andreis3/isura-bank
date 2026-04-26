package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

type ObservabilityTransactionRepo struct {
	repo   transaction.Repository
	metric application.Metrics
	tracer application.Tracer
}

func NewObservabilityTransactionRepo(repo transaction.Repository, metric application.Metrics, tracer application.Tracer) *ObservabilityTransactionRepo {
	return &ObservabilityTransactionRepo{
		repo:   repo,
		metric: metric,
		tracer: tracer,
	}
}

func (r *ObservabilityTransactionRepo) Save(ctx context.Context, data *transaction.Transaction) error {
	ctx, span := r.tracer.Start(ctx, "TransactionRepository.Save")
	defer span.End()

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
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityTransactionRepo) FindByID(ctx context.Context, transactionID transaction.TransactionID) (*transaction.Transaction, error) {
	ctx, span := r.tracer.Start(ctx, "TransactionRepository.FindByID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"find_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	transactionResponse, err := r.repo.FindByID(ctx, transactionID)

	if err != nil {

		return nil, err
	}

	return transactionResponse, nil
}

func (r *ObservabilityTransactionRepo) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*transaction.Transaction, error) {
	ctx, span := r.tracer.Start(ctx, "TransactionRepository.FindByIdempotencyKey")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"transactions",
			"find_by_idempotency_key",
			float64(time.Since(start).Milliseconds()))
	}()

	transactionResponse, err := r.repo.FindByIdempotencyKey(ctx, idempotencyKey)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return transactionResponse, nil
}

func (r *ObservabilityTransactionRepo) ExistsByIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "TransactionRepository.ExistsByIdempotencyKey")
	defer span.End()

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
		span.RecordError(err)
		return false, err
	}

	return exists, nil
}
