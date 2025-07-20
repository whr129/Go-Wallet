CREATE TABLE IF NOT EXISTS users (
    "id" BIGSERIAL PRIMARY KEY,
    "user_name" VARCHAR(255) NOT NULL UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "hash_password" VARCHAR(255) NOT NULL,
    "role" VARCHAR(255) NOT NULL DEFAULT 'user',
    "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_email ON users (email);