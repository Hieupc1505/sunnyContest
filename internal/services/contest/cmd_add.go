package contest

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/pkg/errsx"
)

func newAddContest(userId, subjectId int64, numQuestion int32, timeExam int32) (*db.CreateContestParams, error) {
	var errs errsx.Map
	time, err := newTimeExam(timeExam)
	if err != nil {
		errs.Set("time_exam", err)
	}
	numq, err := newNumQuestion(numQuestion)
	if err != nil {
		errs.Set("num_question", err)
	}
	if len(errs) > 0 {
		return nil, errs.ToError()
	}

	contest := NewAddContest(userId, subjectId, numq, time)

	return contest, nil

}

func (s *Service) Add(ctx *gin.Context, body AddAndUpdateParams, total int64) (*db.SfContest, error) {
	userId, err := contextutil.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if int64(body.NumQuestion) > total {
		return nil, ErrInvalidContestNumberQuestion
	}

	data, err := newAddContest(userId, body.SubjectID, body.NumQuestion, body.TimeExam)
	if err != nil {
		return nil, err
	}
	contest, err := s.repo.CreateContest(ctx, *data)
	if err != nil {
		return nil, err
	}
	return &contest, nil
}
