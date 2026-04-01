# Isura Bank

> **Isura** — *treasure*, in Yoruba.

A digital banking platform built with microservices architecture, Domain-Driven Design, and Go. The project implements the core domains of a modern bank — core banking, payments, cards, credit, investments, and platform services — with a focus on financial consistency, traceability, and regulatory compliance (BACEN).

---

## Repositories

| Service | Domain | Priority | Status |
|---|---|---|---|
| `isura-ledger-ms` | Double-entry accounting engine | P0 | 🟡 In development |
| `isura-account-ms` | Account lifecycle management | P0 | 🔴 Planned |
| `isura-auth-ms` | Authentication and authorization | P0 | 🔴 Planned |
| `isura-pix-ms` | BACEN SPI/DICT integration | P0 | 🔴 Planned |
| `isura-card-ms` | Card lifecycle management | P0 | 🔴 Planned |
| `isura-card-auth-ms` | Purchase authorization (<100ms) | P0 | 🔴 Planned |
| `isura-antifraud-ms` | Real-time fraud analysis | P0 | 🔴 Planned |
| `isura-transfer-ms` | Wire transfers (TED/TEF) | P1 | 🔴 Planned |
| `isura-boleto-ms` | Boleto issuance and settlement | P1 | 🔴 Planned |
| `isura-card-billing-ms` | Credit card billing | P1 | 🔴 Planned |
| `isura-credit-engine-ms` | Credit decision engine | P1 | 🔴 Planned |
| `isura-loan-ms` | Loan management | P1 | 🔴 Planned |
| `isura-notification-ms` | Notification hub | P1 | 🔴 Planned |
| `isura-reconciliation-ms` | Automated reconciliation | P1 | 🔴 Planned |
| `isura-audit-ms` | Immutable audit log | P1 | 🔴 Planned |
| `isura-invest-ms` | Investment platform | P2 | 🔴 Planned |
| `isura-bank-arch-ui` | Architecture visualizer | — | 🟡 In development |

---

## Architecture

### Technical decisions

**Language:** Go 1.26.1 — chosen for its performance, native concurrency via goroutines and channels, static binaries, and adoption by reference digital banks such as Monzo.

**Inter-service communication:**
- Synchronous via **gRPC + Protobuf** when the caller requires a response to proceed (e.g. `pix-svc` → `ledger-svc` for debit)
- Asynchronous via **Apache Kafka** for eventual consistency scenarios (e.g. `ledger-svc` → `notif-svc` to notify the customer)

**Persistence:** PostgreSQL 16 with *database-per-service* — each microservice owns its schema. No service reads directly from another service's database.

**Distributed consistency:**
- **Transactional Outbox** to guarantee event delivery to Kafka without dual writes
- **Saga with orchestration** for long-running flows such as Pix transfers and loans
- **Idempotency** via `idempotency_key` with `UNIQUE CONSTRAINT` in PostgreSQL

**Shared infrastructure:**

| Component | Role |
|---|---|
| PostgreSQL 16 | Per-service persistence |
| Redis 7 | Cache, rate limiting, feature store |
| Apache Kafka | Event streaming, audit trail, CDC |
| Kubernetes (EKS) | Orchestration and auto-scaling |
| OpenTelemetry | Distributed traces, metrics, and logs |
| Vault / HSM | Secrets and cryptographic keys |

### Dependency graph

```
                    ┌─────────────┐
                    │ API Gateway │
                    └──────┬──────┘
                           │
        ┌──────────┬───────┼───────┬──────────┐
        ▼          ▼       ▼       ▼          ▼
  ┌──────────┐ ┌──────┐ ┌─────┐ ┌──────┐ ┌───────┐
  │ account  │ │ pix  │ │card │ │ loan │ │invest │
  └────┬─────┘ └──┬───┘ └──┬──┘ └──┬───┘ └──┬────┘
       │          │        │       │         │
       └──────────┴────────┴───────┴─────────┘
                           │
                ╔══════════╧══════════╗
                ║    LEDGER SERVICE   ║
                ║  (source of truth)  ║
                ╚═════════════════════╝
```

`ledger-svc` is the leaf in the dependency tree — it has no runtime dependency on any other domain service.

---


### Tech stack

| Layer | Technology                          |
|---|-------------------------------------|
| Transport | gRPC + Protobuf                     |
| Application | Go 1.26.1                           |
| Persistence | PostgreSQL 16 + pgx/v5              |
| Migrations | golang-migrate                      |
| Events | Apache Kafka (Transactional Outbox) |
| Observability | OpenTelemetry                       |
| Container | Docker + Kubernetes                 |

---

## Design principles

**Database per service** — each service owns its database. No service reads directly from another service's database.

**Interface on the consumer side** — Go interfaces are defined in the package that consumes them, not in the package that implements them.

**Errors as first-class citizens** — exported sentinel errors in all domain packages. No anonymous `errors.New` at decision points.

**Explicit over convenient** — Go favors explicitness. No magic, no hidden annotations, no global state.

**No circular dependencies between aggregates** — aggregates reference each other only by ID, never by the full type.

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
