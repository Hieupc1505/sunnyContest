package subject_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
	GetAllSubjects(ctx context.Context, arg db.GetAllSubjectsParams) ([]db.SfSubject, error)
	GetSubjectByID(ctx context.Context, id int64) (db.SfSubject, error)
}

// Writer NewReader returns
type Writer interface {
	AddSubject(ctx context.Context, arg db.AddSubjectParams) (db.AddSubjectRow, error)
	DeleteSubject(ctx context.Context, id int64) error
	UpdateSubject(ctx context.Context, arg db.UpdateSubjectParams) (db.UpdateSubjectRow, error)
}

// ReadWriter NewWriter creates a new Writer
type ReadWriter interface {
	Reader
	Writer
}

// NewSubjectRepo create a new user repo
func NewSubjectRepo(conn *pgxpool.Pool) ReadWriter {
	return db.New(conn)
}
