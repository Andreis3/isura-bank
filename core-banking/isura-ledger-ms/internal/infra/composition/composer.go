package composition

import (
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/observability"
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
	createAccount := command.NewCreateAccount(accountRepo, c.deps.Log, c.deps.Tracer)

	// handlers
	createAccountHandler := handler.NewCreateAccountHandler(createAccount, c.deps.Log, c.deps.Tracer)

	// server
	ledgerServer := grpcTransport.NewLedgerServer(createAccountHandler)

	// server
	return server.NewGRPCServer(c.deps, ledgerServer)
}

func (c *Composer) buildAccountRepo() account.Repository {
	return observability.NewObservabilityAccountRepo(
		repository.NewAccountRepository(c.deps.Pg),
		c.deps.Prom,
		c.deps.Tracer,
	)
}

func (c *Composer) buildTransactionRepo() transaction.Repository {
	return observability.NewObservabilityTransactionRepo(
		repository.NewTransactionRepository(c.deps.Pg),
		c.deps.Prom,
		c.deps.Tracer,
	)
}

func (c *Composer) buildOutboxRepo() outbox.Repository {
	return observability.NewObservabilityOutboxRepo(
		repository.NewOutBoxRepository(c.deps.Pg),
		c.deps.Prom,
		c.deps.Tracer,
	)
}
