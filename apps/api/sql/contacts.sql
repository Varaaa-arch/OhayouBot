-- name: CreateContact :one
INSERT INTO contacts (name, phone)
VALUES ($1, $2)
RETURNING *;

-- name: GetContacts :many
SELECT * FROM contacts
ORDER BY created_at DESC;
