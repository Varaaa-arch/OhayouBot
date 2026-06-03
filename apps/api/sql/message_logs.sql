-- name: CreateMessageLog :one
INSERT INTO message_logs (contact_id, message, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMessageLogsByContact :many
SELECT * FROM message_logs
WHERE contact_id = $1
ORDER BY sent_at DESC;
