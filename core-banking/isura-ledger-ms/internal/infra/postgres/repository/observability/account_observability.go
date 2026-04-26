package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

type ObservabilityAccountRepo struct {
	repo   account.Repository
	metric application.Metrics
	tracer application.Tracer
}

func NewObservabilityAccountRepo(repo account.Repository, metric application.Metrics, tracer application.Tracer) *ObservabilityAccountRepo {
	return &ObservabilityAccountRepo{
		repo:   repo,
		metric: metric,
		tracer: tracer,
	}
}

func (r *ObservabilityAccountRepo) Save(ctx context.Context, account *account.Account) error {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.Save")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"save",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.Save(ctx, account)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityAccountRepo) FindByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindByID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	accountResponse, err := r.repo.FindByID(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return accountResponse, nil
}

func (r *ObservabilityAccountRepo) FindByExternalID(ctx context.Context, externalID string) (*account.Account, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindByExternalID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_external_id",
			float64(time.Since(start).Milliseconds()))
	}()

	accountResponse, err := r.repo.FindByExternalID(ctx, externalID)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return accountResponse, nil
}

func (r *ObservabilityAccountRepo) UpdateBalance(ctx context.Context, accountID account.AccountID, balance money.Money) error {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.UpdateBalance")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"update_balance",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.UpdateBalance(ctx, accountID, balance)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityAccountRepo) FindBalanceByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindBalanceByID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_balance_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	balance, err := r.repo.FindBalanceByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		return money.Money{}, err
	}

	return balance, nil
}

func (r *ObservabilityAccountRepo) FindBalanceForUpdateByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindBalanceForUpdateByID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_balance_for_update_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	balance, err := r.repo.FindBalanceForUpdateByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		return money.Money{}, err
	}

	return balance, nil
}
