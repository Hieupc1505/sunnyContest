-- name: AddProfile :one
INSERT INTO sf_profile (
    user_id, nickname, avatar
) VALUES (
             $1, $2, $3
         ) RETURNING nickname, avatar;