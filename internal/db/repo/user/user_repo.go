package user_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
}

// Writer NewReader returns
type Writer interface {
	Add(ctx context.Context, arg db.AddParams) (db.User, error)
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
