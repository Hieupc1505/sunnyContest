-- name: Add :one
INSERT INTO sf_user (
    username, password
) VALUES (
$1, $2
) RETURNING id, username, role, status, created_time;

-- name: CheckUsernameExists :one
SELECT EXISTS (
    SELECT 1 FROM sf_user WHERE username = $1
);

-- name: GetUserByUsername :one
SELECT
    u.id,
    u.username,
    u.password,
    u.role,
    u.status,
    u.created_time,
    u.updated_time,
    COALESCE(p.nickname, '') AS nickname,
    COALESCE(p.avatar, '') AS avatar
FROM sf_user u
         LEFT JOIN sf_profile p ON u.id = p.user_id
WHERE u.username = $1
    LIMIT 1;

-- name: GetUserByID :one
SELECT
    u.id,
    u.username,
    u.role,
    u.status,
    u.created_time,
    u.updated_time,
    COALESCE(p.nickname, '') AS nickname,
    COALESCE(p.avatar, '') AS avatar
FROM sf_user u
LEFT JOIN sf_profile p ON u.id = p.user_id
WHERE u.id = $1
LIMIT 1;

-- name: UpdateUserToken :exec
UPDATE sf_user
SET token = $2, token_expired = $3, updated_time = now()
WHERE id = $1;

-- name: UpdateUserRole :exec
UPDATE sf_user
SET role = $2, updated_time = now()
WHERE id = $1;
