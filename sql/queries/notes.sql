-- name: CreateNote :one
INSERT INTO notes (id, created_at, updated_at, title, content, user_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetNotes :many
SELECT * FROM notes
ORDER BY created_at ASC;

-- name: GetNoteByID :one
SELECT * FROM notes
WHERE id = $1;

-- name: UpdateNote :one
UPDATE notes
SET title = $1, content = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, title, content;

-- name: DeleteNote :exec
DELETE FROM notes *
WHERE id = $1;