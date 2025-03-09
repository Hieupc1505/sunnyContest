package contest

import (
	"context"
	"errors"
	repo "go-rest-api-boilerplate/internal/db/repo/contest"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

var (
	ErrInvalidContestNumberQuestion = errors.New("invalid_data.contest.number_question")
)

type Service struct {
	repo repo.ReadWriter
	Db   db.DBTX
}

func NewService(ctx context.Context, repo repo.ReadWriter) (*Service, error) {
	return &Service{repo: repo}, nil
}
