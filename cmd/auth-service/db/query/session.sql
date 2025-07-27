-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  user_name,
  email,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions
SET
  user_name = COALESCE(sqlc.narg(user_name), user_name),
  email = COALESCE(sqlc.narg(email), email),
  is_blocked = COALESCE(sqlc.narg(is_blocked), is_blocked)
WHERE
  id = sqlc.arg(id)
RETURNING *;
