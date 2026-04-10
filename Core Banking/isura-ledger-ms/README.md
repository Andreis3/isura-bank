# isura-ledger-ms

The double-entry accounting engine of Isura Bank. **Every financial movement goes through here.**

This service is the financial source of truth — it records all debit and credit entries, maintains real-time account balances, guarantees atomic double-entry bookkeeping, and publishes domain events reliably via the Transactional Outbox pattern.

---

## Responsibilities

- Double-entry bookkeeping — debit + credit = 0, always atomic via a single PostgreSQL transaction
- Real-time balance per account
- Transaction idempotency via `idempotency_key` with `UNIQUE CONSTRAINT`
- Hold/Release for fund reservations (card authorizations, pending Pix)
- Reliable event publishing via Transactional Outbox → Kafka
- Accounting representation of accounts (asset, liability, revenue, expense)

---

## Architecture

This service follows **Hexagonal Architecture** with **tactical Domain-Driven Design**, organized around three bounded contexts internally: `transaction`, `account`, and `outbox`.

### Dependency rule

```
transport → application → domain ← infrastructure
```

The `domain` layer imports nothing external. The `infrastructure` layer implements the interfaces defined in `domain`. The `transport` layer knows only about `application` use cases.

### Project structure

```
isura-ledger-ms/
├── cmd/
│   └── server/
│       └── main.go                  # Composition root — wires all dependencies
├── internal/
│   ├── domain/
│   │   ├── money/
│   │   │   └── money.go             # Value Object: Money (int64 cents + Currency)
│   │   ├── account/
│   │   │   ├── account.go           # Entity: Account
│   │   │   └── repository.go        # Repository interface
│   │   ├── transaction/
│   │   │   ├── types.go             # Shared types: EntryID, TransactionID, AccountID
│   │   │   ├── transaction.go       # Aggregate Root: Transaction + state machine
│   │   │   ├── entry.go             # Entity: Entry (one side of double-entry)
│   │   │   └── repository.go        # Repository interface
│   │   └── outbox/
│   │       ├── outbox.go            # Aggregate: OutboxEvent + state machine
│   │       └── repository.go        # Repository interface
│   ├── application/
│   │   ├── uow.go                   # UnitOfWork interface
│   │   ├── event/
│   │   │   └── transaction_created.go  # Domain event payload (outbound contract)
│   │   └── command/
│   │       └── create_transaction.go   # CreateTransaction use case (CQRS write)
│   ├── infrastructure/
│   │   └── postgres/
│   │       ├── database/
│   │       │   └── querier.go       # Querier interface + tx context helpers
│   │       ├── uow/
│   │       │   └── uow.go           # Concrete UnitOfWork (pgx transaction lifecycle)
│   │       ├── model/
│   │       │   ├── transaction_model.go
│   │       │   ├── entry_model.go
│   │       │   ├── account_model.go
│   │       │   └── outbox_model.go
│   │       └── repository/
│   │           ├── resolve_db.go    # resolveDB: picks tx or pool from context
│   │           ├── transaction.go
│   │           ├── account.go
│   │           └── outbox.go
│   └── transport/
│       └── grpc/
│           └── handler.go           # gRPC handler — translates protobuf ↔ command
├── proto/
│   └── ledger.proto                 # gRPC service definition
├── db/
│   └── migrations/                  # SQL migration files (golang-migrate)
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## Domain model

### Money — Value Object

Monetary values are stored as `int64` cents — never `float64`. This eliminates floating-point precision errors in financial calculations.

```
Money { amount int64, currency Currency }

BRL → "BRL" | USD → "USD" | EUR → "EUR"
```

Operations: `Add`, `Subtract`, `Equal`, `IsSufficientBalance`, `IsZero`, `IsNegative`, `IsPositive`, `String` (e.g. `"100.50 BRL"`).

### Transaction — Aggregate Root

The central aggregate. Enforces double-entry invariants and owns its state machine.

```
PENDING → COMPLETED  ✓
PENDING → FAILED     ✓
COMPLETED → any      ✗
FAILED → any         ✗
```

A `Transaction` always contains exactly two `Entry` records — one `DEBIT` and one `CREDIT` with equal amounts. These invariants are enforced by `AddEntry()` before any persistence occurs.

### Entry — Entity

Represents one side of a double-entry ledger. Belongs to a `Transaction` aggregate — never created independently.

```
Entry {
  ID             EntryID
  TransactionID  TransactionID
  AccountID      AccountID
  Direction      DEBIT | CREDIT
  Amount         Money
  IdempotencyKey string
  CreatedAt      time.Time
}
```

### Account — Entity

Accounting representation of a bank account within the ledger. Not the same as the customer-facing account in `isura-account-ms` — the ledger holds only what it needs for bookkeeping.

```
Account {
  ID             AccountID
  ExternalID     string      // ID from isura-account-ms (correlation key)
  AccountingType ASSET | LIABILITY | REVENUE | EXPENSE
  Balance        Money
  CreatedAt      time.Time
  UpdatedAt      time.Time
}
```

### OutboxEvent — Aggregate

Ensures reliable event delivery to Kafka without dual writes. Persisted in the same PostgreSQL transaction as the `Transaction`. A background relay reads `PENDING` events with `SELECT FOR UPDATE SKIP LOCKED` and publishes to Kafka.

```
PENDING → FAILED    ✓  (publish attempt failed)
PENDING → SUCCESS   ✓  (published successfully)
FAILED  → PENDING   ✓  (retry — if Attempts < MaxAttempts)
SUCCESS → any       ✗
```

---

## Key design decisions

### Transactional Outbox

All writes within `CreateTransaction` happen in a single PostgreSQL transaction:

```
BEGIN
  INSERT INTO transactions ...
  INSERT INTO entries ...      (debit)
  INSERT INTO entries ...      (credit)
  UPDATE accounts SET balance  (debit account)
  UPDATE accounts SET balance  (credit account)
  INSERT INTO outbox_events ... (event payload)
COMMIT
```

If the commit fails, nothing is persisted — including the outbox event. If the commit succeeds, the relay will eventually publish the event to Kafka. No dual write, no inconsistency.

### Idempotency

Every transaction carries an `idempotency_key`. A `UNIQUE CONSTRAINT` on the `transactions` table ensures that concurrent retries with the same key result in exactly one committed transaction — even if two requests race past the application-level `ExistsByIdempotencyKey` check.

### Unit of Work

The `UnitOfWork` interface wraps the PostgreSQL transaction lifecycle. All repository calls within a use case's `Execute` method receive a `context.Context` carrying the active `pgx.Tx`. Each repository's `resolveDB` method picks the transaction over the connection pool when present.

```go
return c.uow.WithTransaction(ctx, func(ctxTx context.Context) error {
    // all writes here share the same pgx.Tx
    c.transactionRepo.Save(ctxTx, tx)
    c.accountRepo.UpdateBalance(ctxTx, ...)
    c.outboxRepo.Save(ctxTx, event)
    return nil
})
```

### CQRS

Write operations live in `application/command/`, read operations in `application/query/`. Commands return only `error`. Queries return `(Result, error)`.

---

## Tech stack

| Layer | Technology                          |
|---|-------------------------------------|
| Transport | gRPC + Protobuf                     |
| Language | Go 1.26.1                           |
| Persistence | PostgreSQL 16 + pgx/v5              |
| Migrations | golang-migrate                      |
| Events | Apache Kafka (Transactional Outbox) |
| Observability | OpenTelemetry                       |
| Container | Docker + Kubernetes                 |

### Key dependencies

```
github.com/jackc/pgx/v5          # PostgreSQL driver
github.com/google/uuid           # UUID generation
github.com/golang-migrate/migrate # Database migrations
google.golang.org/grpc           # gRPC server
google.golang.org/protobuf       # Protobuf serialization
go.opentelemetry.io/otel         # Observability
```

---

## Running locally

### Prerequisites

- Go 1.26.1+
- Docker and Docker Compose
- `golang-migrate` CLI

### Setup

```bash
# clone
git clone https://github.com/andreis3/isura-ledger-ms
cd isura-ledger-ms

# start PostgreSQL and Kafka
docker compose up -d

# run migrations
make migrate-up

# start the service
go run ./cmd/server/main.go
```

### Makefile commands

```bash
make migrate-up      # apply all pending migrations
make migrate-down    # rollback last migration
make test            # run unit tests
make test-int        # run integration tests
make proto           # regenerate protobuf files
make lint            # run golangci-lint
make build           # build binary
```

---

## Environment variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=isura_ledger
DB_USER=postgres
DB_PASSWORD=postgres
DB_MAX_CONNS=20
DB_MIN_CONNS=5

# gRPC
GRPC_PORT=50051

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_LEDGER_EVENTS=ledger.events

# Outbox relay
OUTBOX_RELAY_INTERVAL_MS=500
OUTBOX_RELAY_BATCH_SIZE=50
OUTBOX_MAX_ATTEMPTS=3

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=isura-ledger-ms
```

---

## Database schema

### tables

| Table | Description |
|---|---|
| `transactions` | Aggregate root — one record per transaction |
| `entries` | Double-entry records — always two per transaction |
| `accounts` | Accounting representation of accounts |
| `outbox_events` | Transactional outbox — pending Kafka events |

### Key constraints

```sql
-- idempotency guarantee
UNIQUE (idempotency_key) ON transactions

-- double-entry integrity
FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON entries

-- outbox ordering
INDEX (status, created_at) ON outbox_events
```

---

## gRPC API

```protobuf
service LedgerService {
  rpc CreateTransaction (CreateTransactionRequest) returns (CreateTransactionResponse);
  rpc GetBalance        (GetBalanceRequest)         returns (GetBalanceResponse);
  rpc GetTransaction    (GetTransactionRequest)      returns (GetTransactionResponse);
  rpc CreateAccount     (CreateAccountRequest)       returns (CreateAccountResponse);
}
```

---

## Testing

```bash
# unit tests (domain logic — no external dependencies)
go test ./internal/domain/...

# integration tests (requires Docker)
go test ./internal/infrastructure/... -tags=integration

# all tests with race detector
go test -race ./...
```

---

## References

- *Implementing Domain-Driven Design* — Vaughn Vernon
- *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi
- *Concurrency in Go* — Katherine Cox-Buday
- *The Go Programming Language* — Donovan & Kernighan
- [pgx/v5 documentation](https://github.com/jackc/pgx)
- [gRPC Go documentation](https://grpc.io/docs/languages/go/)

---

## License

MIT — see [LICENSE](./LICENSE) for details.
