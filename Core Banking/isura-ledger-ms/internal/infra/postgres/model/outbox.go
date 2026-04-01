package model

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
)

type Outbox struct {
	ID            pgtype.Text
	AggregateID   pgtype.Text
	AggregateType pgtype.Text
	EventType     pgtype.Text
	Payload       []byte
	Status        pgtype.Text
	Attempts      pgtype.Int2
	LastAttemptAt pgtype.Timestamptz
	CreatedAt     pgtype.Timestamptz
	PublishedAt   pgtype.Timestamptz
}

func ToOutboxModel(domain *outbox.Outbox) Outbox {
	return Outbox{
		ID: pgtype.Text{
			String: string(domain.ID),
			Valid:  true,
		},
		Status: pgtype.Text{
			String: string(domain.Status),
			Valid:  true,
		},
		AggregateID: pgtype.Text{
			String: string(domain.AggregateID),
			Valid:  true,
		},
		AggregateType: pgtype.Text{
			String: string(domain.AggregateType),
			Valid:  true,
		},
		Attempts: pgtype.Int2{
			Int16: int16(domain.Attempts),
			Valid: true,
		},
		EventType: pgtype.Text{
			String: string(domain.EventType),
			Valid:  true,
		},
		Payload:       domain.Payload,
		LastAttemptAt: database.ToTimestamptz(domain.LastAttemptAt),
		CreatedAt: pgtype.Timestamptz{
			Time:  domain.CreatedAt,
			Valid: true,
		},
		PublishedAt: database.ToTimestamptz(domain.PublishedAt),
	}
}

func ToOutboxDomain(model Outbox) *outbox.Outbox {
	return &outbox.Outbox{
		ID:            outbox.OutboxID(model.ID.String),
		EventType:     outbox.EventType(model.EventType.String),
		Attempts:      int(model.Attempts.Int16),
		AggregateType: outbox.AggregateType(model.AggregateType.String),
		AggregateID:   string(model.AggregateID.String),
		Status:        outbox.StatusOutbox(model.Status.String),
		Payload:       model.Payload,
		CreatedAt:     model.CreatedAt.Time,
		LastAttemptAt: database.ToTimePtr(model.LastAttemptAt),
		PublishedAt:   database.ToTimePtr(model.PublishedAt),
	}
}
