-- name: CreateSession :one
INSERT INTO whatsapp_sessions (session_data)
VALUES ($1)
RETURNING *;

-- name: GetLatestSession :one
SELECT * FROM whatsapp_sessions
ORDER BY created_at DESC
LIMIT 1;
