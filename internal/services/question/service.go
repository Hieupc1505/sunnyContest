package question

import (
	"context"
	repo "go-rest-api-boilerplate/internal/db/repo/question"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

type Service struct {
	repo repo.ReadWriter
	Db   db.DBTX
}

func NewService(ctx context.Context, repo repo.ReadWriter) (*Service, error) {
	return &Service{repo: repo}, nil
}
