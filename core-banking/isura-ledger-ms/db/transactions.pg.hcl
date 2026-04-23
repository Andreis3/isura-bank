table "transactions" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "idempotency_key" {
    type = varchar(50)
    null = false
  }

  column "status" {
    type = varchar(20)
    null = false
  }

  column "created_at" {
    type = timestamptz
    null = false
  }

  column "updated_at" {
    type = timestamptz
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_transactions_idempotency_key" {
    columns = [column.idempotency_key]
    unique  = true
  }

  index "idx_transactions_status" {
    columns = [column.status]
  }
}