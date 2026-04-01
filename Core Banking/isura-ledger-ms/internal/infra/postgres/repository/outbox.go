package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/andreis3/isura-ledger-ms/internal/domain/outbox"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/model"
)

type OutBoxRepository struct {
	db database.Querier
}

func NewOutBoxRepository(db database.Querier) *OutBoxRepository {
	return &OutBoxRepository{
		db: db,
	}
}

func (r *OutBoxRepository) Save(ctx context.Context, outbox *outbox.Outbox) error {
	db := resolveDB(ctx, r.db)

	query := `INSERT INTO outboxs (
		id,
		aggregate_id,
		aggregate_type,
		event_type,
        payload,
    	status,
    	attempts,
    	last_attempt_at,
    	created_at,
    	published_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	outboxModel := model.ToOutboxModel(outbox)

	_, err := db.Exec(ctx, query,
		outboxModel.ID,
		outboxModel.AggregateID,
		outboxModel.AggregateType,
		outboxModel.EventType,
		outboxModel.Payload,
		outboxModel.Status,
		outboxModel.Attempts,
		outboxModel.LastAttemptAt,
		outboxModel.CreatedAt,
		outboxModel.PublishedAt,
	)

	return err

}
func (r *OutBoxRepository) FindAllByStatusForUpdateSkipLocked(ctx context.Context, status outbox.StatusOutbox, limit int) ([]*outbox.Outbox, error) {
	db := resolveDB(ctx, r.db)

	query := `
	SELECT 
		id,
		aggregate_id,
		aggregate_type,
		event_type,
		payload,
		status,
		attempts,
		last_attempt_at,
		created_at,
		published_at
	FROM outboxs
	WHERE status = $1
	LIMIT $2
	FOR UPDATE SKIP LOCKED
	`

	rows, err := db.Query(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var outboxes []*outbox.Outbox
	for rows.Next() {
		var outboxModel model.Outbox
		if err := rows.Scan(
			&outboxModel.ID,
			&outboxModel.AggregateID,
			&outboxModel.AggregateType,
			&outboxModel.EventType,
			&outboxModel.Payload,
			&outboxModel.Status,
			&outboxModel.Attempts,
			&outboxModel.LastAttemptAt,
			&outboxModel.CreatedAt,
			&outboxModel.PublishedAt,
		); err != nil {
			return nil, err
		}
		outbox := model.ToOutboxDomain(outboxModel)
		outboxes = append(outboxes, outbox)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return outboxes, nil
}
func (r *OutBoxRepository) UpdateOutboxData(ctx context.Context, outboxID outbox.OutboxID, data outbox.UpdateOutboxData) error {
	db := resolveDB(ctx, r.db)

	query := `
	UPDATE outboxs
	SET status = $1, attempts = $2, last_attempt_at = $3, published_at = $4
	WHERE id = $5
	`

	_, err := db.Exec(ctx, query,
		pgtype.Text{String: string(data.Status), Valid: true},
		pgtype.Int2{Int16: int16(data.Attempts), Valid: true},
		database.ToTimestamptz(data.LastAttemptAt),
		database.ToTimestamptz(data.PublishedAt),
		pgtype.Text{String: string(outboxID), Valid: true},
	)

	return err
}
