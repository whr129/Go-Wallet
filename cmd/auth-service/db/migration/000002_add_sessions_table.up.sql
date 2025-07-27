CREATE TABLE "sessions" (
  "id" bigint PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "user_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT FALSE,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);