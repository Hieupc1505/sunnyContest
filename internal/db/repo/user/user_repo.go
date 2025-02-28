package user_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error)
}

// Writer NewReader returns
type Writer interface {
	Add(ctx context.Context, arg db.AddParams) (db.AddRow, error)
	UpdateUserToken(ctx context.Context, arg db.UpdateUserTokenParams) error
	AddProfile(ctx context.Context, arg db.AddProfileParams) (db.AddProfileRow, error)
	UpdateUserRole(ctx context.Context, arg db.UpdateUserRoleParams) error
}

// ReadWriter NewWriter creates a new Writer
type ReadWriter interface {
	Reader
	Writer
}

// NewUserRepo create a new user repo
func NewUserRepo(conn *pgxpool.Pool) ReadWriter {
	return db.New(conn)
}
