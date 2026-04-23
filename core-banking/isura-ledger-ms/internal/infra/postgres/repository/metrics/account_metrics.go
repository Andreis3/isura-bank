package metrics

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

type MetricsAccountRepo struct {
	repo   account.Repository
	metric application.Metrics
}

func NewMetricsAccountRepo(repo account.Repository, metric application.Metrics) *MetricsAccountRepo {
	return &MetricsAccountRepo{
		repo:   repo,
		metric: metric,
	}
}

func (r *MetricsAccountRepo) Save(ctx context.Context, account *account.Account) error {
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
		return err
	}

	return nil
}

func (r *MetricsAccountRepo) FindByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	account, err := r.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *MetricsAccountRepo) FindByExternalID(ctx context.Context, externalID string) (*account.Account, error) {
	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_external_id",
			float64(time.Since(start).Milliseconds()))
	}()

	account, err := r.repo.FindByExternalID(ctx, externalID)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *MetricsAccountRepo) UpdateBalance(ctx context.Context, accountID account.AccountID, balance money.Money) error {
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
		return err
	}

	return nil
}

func (r *MetricsAccountRepo) FindBalanceByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
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
		return money.Money{}, err
	}

	return balance, nil
}

func (r *MetricsAccountRepo) FindBalanceForUpdateByID(ctx context.Context, accountID account.AccountID) (money.Money, error) {
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
		return money.Money{}, err
	}

	return balance, nil
}
