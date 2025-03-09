package user_contest_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
	GetUserContest(ctx context.Context, arg db.GetUserContestParams) (db.SfUserContest, error)
	GetUserContestsByContestID(ctx context.Context, contestID int64) ([]db.GetUserContestsByContestIDRow, error)
	GetUsersInContest(ctx context.Context, contestID int64) ([]db.GetUsersInContestRow, error)
	GetUserContestsJoined(ctx context.Context, userID int64) ([]db.GetUserContestsJoinedRow, error)
}

// Writer NewReader returns
type Writer interface {
	AddUserContest(ctx context.Context, arg db.AddUserContestParams) (db.SfUserContest, error)
	UpdateExamAndResult(ctx context.Context, arg db.UpdateExamAndResultParams) error
}

// ReadWriter NewWriter creates a new Writer
type ReadWriter interface {
	Reader
	Writer
}

// NewSubjectRepo create a new user repo
func NewUserContestService(conn *pgxpool.Pool) ReadWriter {
	return db.New(conn)
}
