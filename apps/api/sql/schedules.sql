-- name: CreateSchedule :one
INSERT INTO schedules (contact_id, message, time, enabled)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSchedules :many
SELECT * FROM schedules
ORDER BY created_at DESC;
