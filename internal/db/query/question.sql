-- name: AddQuestion :one
INSERT INTO sf_question (
    subject_id, user_id, level, question, question_type, question_image,
    answers, answer_type, state, updated_time
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, now()
         )
RETURNING *;

-- name: UpdateQuestion :one
UPDATE sf_question
SET
    subject_id = $2,
    level = $3,
    question = $4,
    question_type = $5,
    question_image = $6,
    answers = $7,
    answer_type = $8,
    user_id = COALESCE(NULLIF(@user_id, 0), user_id),
    state = COALESCE(NULLIF(@state, 0), state),
    updated_time = now()
WHERE id = $1
RETURNING *;

-- name: GetTotalQuestion :one
SELECT COUNT(*) FROM sf_question WHERE subject_id = $1;

-- name: GetQuestionBySubjectID :many
SELECT * FROM sf_question WHERE subject_id = $1 ORDER BY created_time DESC LIMIT $2 OFFSET $3;

-- name: GetQuestionByID :one
SELECT * FROM sf_question WHERE id = $1;

-- name: DeleteQuestion :exec
DELETE FROM sf_question WHERE id = $1;

