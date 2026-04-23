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

	accountRepo := c.buildAccountRepo()

	// use cases
	createAccount := command.NewCreateAccount(accountRepo, c.deps.Log)

	// handlers
	createAccountHandler := handler.NewCreateAccountHandler(createAccount, c.deps.Log)

	// server
	ledgerServer := grpcTransport.NewLedgerServer(createAccountHandler)

	// server
	return server.NewGRPCServer(c.deps, ledgerServer)
}

func (c *Composer) buildAccountRepo() account.Repository {
	return metrics.NewMetricsAccountRepo(
		repository.NewAccountRepository(c.deps.Pg),
		c.deps.Prom,
	)
}

func (c *Composer) buildTransactionRepo() transaction.Repository {
	return metrics.NewMetricsTransactionRepo(
		repository.NewTransactionRepository(c.deps.Pg),
		c.deps.Prom,
	)
}

func (c *Composer) buildOutboxRepo() outbox.Repository {
	return metrics.NewMetricsOutboxRepo(
		repository.NewOutBoxRepository(c.deps.Pg),
		c.deps.Prom,
	)
}
