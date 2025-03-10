package account

import (
	"context"
	"errors"
	repo "go-rest-api-boilerplate/internal/db/repo/user"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

var (
	ErrUsernameInvalidLen = errors.New("invalid username length")
	ErrUsernameEmpty      = errors.New("empty username")
	ErrPasswordInvalidLen = errors.New("invalid password length")
)

type Service struct {
	repo repo.ReadWriter
	Db   db.DBTX
}

func NewService(ctx context.Context, repo repo.ReadWriter) (*Service, error) {
	return &Service{repo: repo}, nil
}
