package subject

import (
	"context"
	repo "go-rest-api-boilerplate/internal/db/repo/subject"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

const (
	StateDefault = 1
)

type Service struct {
	repo repo.ReadWriter
	Db   db.DBTX
}

func NewService(ctx context.Context, repo repo.ReadWriter) (*Service, error) {
	return &Service{repo: repo}, nil
}
