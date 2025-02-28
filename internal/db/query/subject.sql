-- name: AddSubject :one
INSERT INTO sf_subject (
    user_id, name, description, tags, state, created_time
) VALUES (
             $1, $2, $3, $4, $5, now()
         )
    RETURNING id, name, description, tags, state, created_time;

-- name: UpdateSubject :one
UPDATE sf_subject
SET
    name = $2,
    description = $3,
    tags = $4,
    updated_time = now()
WHERE id = $1
RETURNING id, user_id, name, description, tags, state, created_time;

-- name: GetSubjectByID :one
SELECT * FROM sf_subject WHERE id = $1 LIMIT 1;

-- name: DeleteSubject :exec
DELETE FROM sf_subject WHERE id = $1;

-- name: GetAllSubjects :many
SELECT * FROM sf_subject
ORDER BY created_time DESC
    LIMIT $1 OFFSET $2;

