package question

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/pkg/errsx"
	"go-rest-api-boilerplate/types"
	"log/slog"
)

func newQuestion(subjectID int, q, answerType, level, questionImage, questionType string, answers []types.AnswerItem) (*QuestionParams, error) {
	var errs errsx.Map
	nq, err := NewQuestion(q)
	if err != nil {
		errs.Set("question", err)
	}
	l, err := NewLevel(level)
	if err != nil {
		errs.Set("level", err)
	}
	anst, err := NewAnswerType(answerType)
	if err != nil {
		errs.Set("answer_type", err)
	}
	img, err := NewImage(questionImage)
	if err != nil {
		errs.Set("question_image", err)
	}
	qtype, err := NewQuestionType(questionType)
	if err != nil {
		errs.Set("question_type", err)
	}
	ans, err := NewAnswers(answers)
	if err != nil {
		errs.Set("answers", ans)
	}
	if len(errs) > 0 {
		return nil, errs
	}

	params := NewQuestionParams(subjectID, l, nq, qtype, img, anst, *ans)
	return params, nil

}

func (s *Service) AddQuestion(ctx *gin.Context, data QuestionParams) (*db.SfQuestion, error) {
	params, err := newQuestion(data.SubjectID, data.Question, data.AnswerType, data.Level, data.QuestionImage, data.QuestionType, data.Answers)
	if err != nil {
		return nil, err
	}
	userID, err := contextutil.GetUser(ctx)
	if err != nil {
		return nil, app.ErrPermissionDenied
	}
	d := ConvertToAddQuestionParams(*params, userID, stateQuestionDefault)
	question, err := s.repo.AddQuestion(ctx, d)
	if err != nil {
		slog.Error("Error adding question: %v", err)
		return nil, app.ErrInternalServerError
	}
	return &question, nil
}
