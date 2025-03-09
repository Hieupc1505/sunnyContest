-- name: AddUserContest :one
INSERT INTO sf_user_contest (contest_id, user_id, questions)
VALUES (
   $1, -- contest_id
   $2, -- user_id
   $3 -- questions (JSONB)
)
RETURNING *;

-- name: GetUserContest :one
SELECT *
FROM sf_user_contest
WHERE contest_id = $1 AND user_id = $2;

-- name: UpdateExamAndResult :exec
UPDATE sf_user_contest
SET exam = @Exam::jsonb,
    result = @Result::jsonb,
    updated_time = now()
WHERE contest_id = $1 AND user_id = $2;

-- name: GetUserContestsByContestID :many
SELECT
    s.id,
    s.contest_id,
    s.user_id,
    s.questions,
    s.exam,
    s.result,
    s.created_time,
    s.updated_time,
    p.nickname  -- Lấy nickname từ bảng sf_profile
FROM sf_user_contest s
JOIN sf_profile p ON s.user_id = p.user_id  -- Kết nối bảng sf_user_contest với bảng sf_profile qua user_id
WHERE s.contest_id = $1;

-- name: GetUsersInContest :many
SELECT c.*, s.nickname, s.avatar FROM sf_user_contest c
JOIN sf_profile s ON c.user_id = s.user_id
WHERE c.contest_id = $1;


-- name: GetUserContestsJoined :many
SELECT
    c.id,
    c.time_exam,
    c.num_question,
    s.name AS subject_name
FROM
    sf_user_contest uc
        JOIN
    sf_contest c ON uc.contest_id = c.id
        JOIN
    sf_subject s ON c.subject_id = s.id
WHERE
    uc.user_id = $1;
