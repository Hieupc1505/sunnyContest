// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SfProfile struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type SfSubject struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	Tags        pgtype.Text `json:"tags"`
	State       int32       `json:"state"`
	CreatedTime time.Time   `json:"created_time"`
	UpdatedTime time.Time   `json:"updated_time"`
}

type SfUser struct {
	ID           int64              `json:"id"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	Role         int32              `json:"role"`
	Status       int32              `json:"status"`
	Token        pgtype.Text        `json:"token"`
	TokenExpired pgtype.Timestamptz `json:"token_expired"`
	CreatedTime  time.Time          `json:"created_time"`
	UpdatedTime  time.Time          `json:"updated_time"`
}
