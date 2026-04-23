package transaction

import "context"

type Repository interface {
	Save(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, transactionID TransactionID) (*Transaction, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*Transaction, error)
	ExistsByIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error)
}
