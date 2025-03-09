package contest

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/pkg/errsx"
)

func newUpdateContest(id, userId, subjectId int64, numQuestion int32, timeExam int32) (*db.UpdateContestParams, error) {
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

	contest := NewUpdateContest(id, userId, subjectId, numq, time)

	return contest, nil

}

func (s *Service) Update(ctx *gin.Context, body AddAndUpdateParams, subjectQuestionTotal int64) (*db.SfContest, error) {
	userId, err := contextutil.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if int64(body.NumQuestion) > subjectQuestionTotal {
		return nil, ErrInvalidContestNumberQuestion
	}

	data, err := newUpdateContest(body.ID, userId, body.SubjectID, body.NumQuestion, body.TimeExam)
	if err != nil {
		return nil, err
	}
	contest, err := s.repo.UpdateContest(ctx, *data)
	if err != nil {
		return nil, err
	}
	return &contest, nil
}
