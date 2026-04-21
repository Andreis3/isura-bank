package composition

import (
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/metrics"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
	grpcTransport "github.com/andreis3/isura-ledger-ms/internal/transport/grpc"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/handler"
)

type Composer struct {
	deps *server.BaseDeps
}

func NewComposer(baseDeps *server.BaseDeps) *Composer {
	return &Composer{
		deps: baseDeps,
	}
}

func (c *Composer) GRPCServer() *server.GRPCServer {

	// repositories
	accountRepo := c.buildAccountRepo()
	//transactionRepo := c.buildTransactionRepo()
	//outboxRepo := c.buildOutboxRepo()
	//uow := uow.NewUnitOfWork(c.deps.Postgres)

	// use cases
	createAccount := command.NewCreateAccount(accountRepo, c.deps.Logger)
	//createTransaction := command.NewCreateTransaction(uow, accountRepo, transactionRepo, outboxRepo, c.deps.Logger)

	// handlers
	createAccountHandler := handler.NewCreateAccountHandler(createAccount, c.deps.Logger)
	//createTransactionHandler := grpcTransport.NewCreateTransactionHandler(createTransaction, c.deps.Logger)

	// server
	ledgerServer := grpcTransport.NewLedgerServer(createAccountHandler)

	// server
	return server.NewGRPCServer(c.deps.Config, c.deps.Logger, c.deps.Prometheus, c.deps.Postgres, ledgerServer)
}

func (c *Composer) buildAccountRepo() account.Repository {
	return metrics.NewMetricsAccountRepo(
		repository.NewAccountRepository(c.deps.Postgres),
		c.deps.Prometheus,
	)
}

func (c *Composer) buildTransactionRepo() transaction.Repository {
	return metrics.NewMetricsTransactionRepo(
		repository.NewTransactionRepository(c.deps.Postgres),
		c.deps.Prometheus,
	)
}

func (c *Composer) buildOutboxRepo() outbox.Repository {
	return metrics.NewMetricsOutboxRepo(
		repository.NewOutBoxRepository(c.deps.Postgres),
		c.deps.Prometheus,
	)
}
