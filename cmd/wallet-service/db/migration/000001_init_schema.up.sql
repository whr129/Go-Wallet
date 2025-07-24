CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY,
  "user_id" bigint,
  "email" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_account_user_id ON accounts (user_id);