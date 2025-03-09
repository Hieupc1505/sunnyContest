package contest_repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

// Reader Create a reader interface
type Reader interface {
	GetContestLiveByID(ctx context.Context, id int64) (db.GetContestLiveByIDRow, error)
	GetListLiveContest(ctx context.Context) ([]db.GetListLiveContestRow, error)
	GetMyContestLive(ctx context.Context, userID int64) (db.GetMyContestLiveRow, error)
	GetContestByID(ctx context.Context, id int64) (db.SfContest, error)
	GetUsersInContest(ctx context.Context, contestID int64) ([]db.GetUsersInContestRow, error)
	GetContestDetailByID(ctx context.Context, id int64) (db.GetContestDetailByIDRow, error)
	GetContestsForTeacher(ctx context.Context, arg db.GetContestsForTeacherParams) ([]db.GetContestsForTeacherRow, error)
}

// Writer NewReader returns
type Writer interface {
	StartContest(ctx context.Context, arg db.StartContestParams) error
	StopContest(ctx context.Context, id int64) error
	UpdateContest(ctx context.Context, arg db.UpdateContestParams) (db.SfContest, error)
	CreateContest(ctx context.Context, arg db.CreateContestParams) (db.SfContest, error)
	UpdateStateContest(ctx context.Context, arg db.UpdateStateContestParams) error
	UpdateContestStateAndQuestions(ctx context.Context, arg db.UpdateContestStateAndQuestionsParams) error
}

// ReadWriter NewWriter creates a new Writer
type ReadWriter interface {
	Reader
	Writer
}

// NewSubjectRepo create a new user repo
func NewContestService(conn *pgxpool.Pool) ReadWriter {
	return db.New(conn)
}
