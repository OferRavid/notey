-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    null
)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;

-- name: ClearRevokedTokens :exec
DELETE FROM refresh_tokens *
WHERE revoked_at is not null;

-- name: DeleteTokensByUserID :exec
DELETE FROM refresh_tokens *
WHERE user_id = $1;

-- name: CheckRecordExists :one
SELECT COUNT(*) as count
FROM refresh_tokens
WHERE user_id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
