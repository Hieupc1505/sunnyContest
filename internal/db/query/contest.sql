-- name: CreateContest :one
INSERT INTO sf_contest (user_id, subject_id, num_question, time_exam, time_start_exam, state, questions, created_time, updated_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
    RETURNING *;

-- name: UpdateContest :one
UPDATE sf_contest
SET user_id = $2, subject_id = $3, num_question = $4, time_exam = $5, time_start_exam = $6, state = $7, questions = $8, updated_time = now()
WHERE id = $1
RETURNING *;

-- name: UpdateContestQuestions :exec
UPDATE sf_contest
SET questions = $2
WHERE id = $1;

-- name: StartContest :exec
UPDATE sf_contest
SET state = 'RUNNING', time_start_exam = $2
WHERE id = $1;

-- name: StopContest :exec
UPDATE sf_contest
SET state = 'FINISHED'
WHERE id = $1;

-- name: UpdateStateContest :exec
UPDATE sf_contest
SET state = $2, updated_time = now()
WHERE id = $1;

-- name: UpdateContestStateAndQuestions :exec
UPDATE sf_contest
SET state = $2, questions = $3, updated_time = now()
WHERE id = $1;

-- name: GetMyContestLive :one
SELECT c.*, s.name AS subject_name
FROM sf_contest c
JOIN sf_subject s ON c.subject_id = s.id
WHERE c.user_id = $1 AND c.state IN ('IDLE', 'RUNNING', 'WAITING');

-- name: GetListLiveContest :many
SELECT c.id, c.time_exam, c.num_question, c.state, s.name AS subject_name  FROM sf_contest c
JOIN sf_subject s ON c.subject_id = s.id
WHERE c.state IN ('RUNNING', 'WAITING');

-- name: GetContestLiveByID :one
SELECT c.id, c.time_exam, c.num_question, c.state, c.user_id, c.questions, c.subject_id, s.name AS subject_name  FROM sf_contest c
JOIN sf_subject s ON c.subject_id = s.id
WHERE c.id = $1 AND c.state IN ('IDLE', 'RUNNING', 'WAITING');

-- name: GetContestDetailByID :one
SELECT c.id, c.time_exam, c.num_question, c.state, c.user_id, c.questions, c.subject_id, s.name AS subject_name  FROM sf_contest c
JOIN sf_subject s ON c.subject_id = s.id
WHERE c.id = $1;

-- name: GetContestByID :one
SELECT * FROM sf_contest
WHERE id = $1;

-- name: GetContestBySubjectID :many
SELECT * FROM sf_contest
WHERE subject_id = $1;

-- name: GetRandomQuestions :many
SELECT *
FROM sf_question
WHERE subject_id = $1
ORDER BY random()
LIMIT $2;

-- name: GetContestsForTeacher :many
SELECT c.id, c.time_exam, c.num_question, c.state, c.user_id, c.questions, c.subject_id, s.name AS subject_name
FROM sf_contest c
JOIN sf_subject s ON c.subject_id = s.id
WHERE c.user_id = $1
ORDER BY c.created_time
DESC LIMIT $2 OFFSET $3;