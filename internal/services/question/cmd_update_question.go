package question

import (
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"log/slog"
)

func (s *Service) UpdateQuestion(ctx *gin.Context, data QuestionParams) (*db.SfQuestion, error) {
	params, err := newQuestion(data.SubjectID, data.Question, data.AnswerType, data.Level, data.QuestionImage, data.QuestionType, data.Answers)
	if err != nil {
		return nil, err
	}
	params.ID = data.ID

	dataConv := ConvertToUpdateQuestionParams(*params, stateQuestionDefault)
	question, err := s.repo.UpdateQuestion(ctx, dataConv)
	if err != nil {
		slog.Error("Error update question", "error", err)
		return nil, app.ErrInternalServerError
	}
	return &question, nil
}
