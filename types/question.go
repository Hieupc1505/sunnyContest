package types

import (
	"time"
)

type Question struct {
	ID            int64        `json:"id"`
	SubjectID     int64        `json:"subject_id"`
	UserID        int64        `json:"user_id"`
	Level         string       `json:"level"`
	Question      string       `json:"question"`
	QuestionType  string       `json:"question_type"`
	QuestionImage string       `json:"question_image"`
	Answers       []AnswerItem `json:"answers"`
	AnswerType    string       `json:"answer_type"`
	State         int32        `json:"state"`
	CreatedTime   time.Time    `json:"created_time"`
	UpdatedTime   time.Time    `json:"updated_time"`
}
