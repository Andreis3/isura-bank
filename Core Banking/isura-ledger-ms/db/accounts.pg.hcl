table "accounts" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "external_id" {
    type = text
    null = false
  }

  column "account_type" {
    type = varchar(20)
    null = false
  }

  column "balance" {
    type = bigint
    null = false
  }

  column "currency" {
    type = varchar(5)
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

  index "idx_accounts_external_id" {
    columns = [column.external_id]
    unique  = true
  }
}