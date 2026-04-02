table "entries" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "idempotency_key" {
    type = varchar(50)
    null = false
  }

  column "direction" {
    type = varchar(6)
    null = false
  }

  column "amount" {
    type = bigint
    null = false
  }

  column "currency" {
    type = varchar(5)
    null = false
  }

  column "account_id" {
    type = text
    null = false
  }

  column "transaction_id" {
    type = text
    null = false
  }

  column "created_at" {
    type = timestamptz
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_account_id" {
    columns = [column.account_id]
    ref_columns = [table.accounts.column.id]
  }

  foreign_key "fk_transactios_id" {
    columns = [column.transaction_id]
    ref_columns = [table.transactions.column.id]
  }

  index "idx_entries_transaction_id" {
    columns = [column.transaction_id]
  }
}