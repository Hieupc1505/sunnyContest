package question

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/types"
)

const (
	stateQuestionDefault = 1
)

type QuestionParams struct {
	ID            int64              `json:"id"`
	SubjectID     int                `json:"subject_id"`
	Level         string             `json:"level"`
	Question      string             `json:"question"`
	QuestionImage string             `json:"question_image"`
	QuestionType  string             `json:"question_type"`
	AnswerType    string             `json:"answer_type"`
	Answers       []types.AnswerItem `json:"answers"`
}

func NewSfQuestion(subjectID int64, level Level, question Question, questionType QuestionType, img Image, answerType AnswerType, answers Answers) *db.SfQuestion {
	return &db.SfQuestion{
		SubjectID:     subjectID,
		Level:         db.LevelQuestion(level.String()), // Nếu LevelQuestion là enum, cần kiểm tra h��p lệ
		Question:      question.String(),
		QuestionType:  questionType.String(),
		QuestionImage: pgtype.Text{String: img.String(), Valid: true},
		AnswerType:    answerType.String(),
		Answers:       answers.String(), // Chuyển []AnswerItem thành string nếu cần
	}
}

func NewQuestionParams(subjectID int, level Level, question Question, questionType QuestionType, img Image, answerType AnswerType, answers Answers) *QuestionParams {
	return &QuestionParams{
		SubjectID:     subjectID,
		Level:         level.String(),
		Question:      question.String(),
		QuestionImage: img.String(),
		QuestionType:  questionType.String(),
		AnswerType:    answerType.String(),
		Answers:       answers,
	}
}

func ConvertToAddQuestionParams(q QuestionParams, userID int64, state int32) db.AddQuestionParams {
	return db.AddQuestionParams{
		SubjectID:     int64(q.SubjectID),
		UserID:        userID,
		Level:         db.LevelQuestion(q.Level), // Nếu LevelQuestion là enum, cần kiểm tra hợp lệ
		Question:      q.Question,
		QuestionType:  q.QuestionType,
		QuestionImage: pgtype.Text{String: q.QuestionImage, Valid: true},
		Answers:       ConvertAnswersToString(q.Answers), // Chuyển []AnswerItem thành string nếu cần
		AnswerType:    q.AnswerType,
		State:         state,
	}
}

func ConvertToUpdateQuestionParams(q QuestionParams, state int32) db.UpdateQuestionParams {
	return db.UpdateQuestionParams{
		ID:            q.ID,
		SubjectID:     int64(q.SubjectID),
		UserID:        nil,
		Level:         db.LevelQuestion(q.Level),
		Question:      q.Question,
		QuestionType:  q.QuestionType,
		QuestionImage: pgtype.Text{String: q.QuestionImage, Valid: true},
		Answers:       ConvertAnswersToString(q.Answers),
		AnswerType:    q.AnswerType,
		State:         nil,
	}
}

func NewRandomQuestionParams(subjectID int64, limit int32) *db.GetRandomQuestionsParams {
	return &db.GetRandomQuestionsParams{
		SubjectID: subjectID,
		Limit:     limit,
	}
}

func NewQuestionFromSfQuestion(question db.SfQuestion) *types.Question {
	return &types.Question{
		ID:            question.ID,
		SubjectID:     question.SubjectID,
		UserID:        question.UserID,
		Level:         string(question.Level),
		Question:      question.Question,
		QuestionType:  question.QuestionType,
		QuestionImage: question.QuestionImage.String,
		Answers:       ConvertAnswersToJson(question.Answers),
		AnswerType:    question.AnswerType,
		State:         question.State,
		CreatedTime:   question.CreatedTime,
		UpdatedTime:   question.UpdatedTime,
	}
}

// Hàm giả định để chuyển đổi danh sách câu trả lời thành string
func ConvertAnswersToString(answers []types.AnswerItem) string {
	bytes, _ := json.Marshal(answers)
	return string(bytes)
}

func ConvertAnswersToJson(answers string) []types.AnswerItem {
	var result []types.AnswerItem
	_ = json.Unmarshal([]byte(answers), &result)
	return result
}
