-- name: CreateUser :one
INSERT INTO users (
    id,
    user_name,
    hash_password,
    email,
    role
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  hash_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  email = COALESCE(sqlc.narg(email), email),
  role = COALESCE(sqlc.narg(role), role),
  is_deleted = COALESCE(sqlc.narg(is_deleted), is_deleted)
WHERE
  id = sqlc.arg(id)
RETURNING *;
