package question_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
	GetQuestionByID(ctx context.Context, id int64) (db.SfQuestion, error)
	GetQuestionBySubjectID(ctx context.Context, arg db.GetQuestionBySubjectIDParams) ([]db.SfQuestion, error)
	GetTotalQuestion(ctx context.Context, subjectID int64) (int64, error)
}

// Writer NewReader returns
type Writer interface {
	AddQuestion(ctx context.Context, arg db.AddQuestionParams) (db.SfQuestion, error)
	DeleteQuestion(ctx context.Context, id int64) error
	UpdateQuestion(ctx context.Context, arg db.UpdateQuestionParams) (db.SfQuestion, error)
}

// ReadWriter NewWriter creates a new Writer
type ReadWriter interface {
	Reader
	Writer
}

// NewSubjectRepo create a new user repo
func NewQuestionService(conn *pgxpool.Pool) ReadWriter {
	return db.New(conn)
}
