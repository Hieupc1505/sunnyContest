package question

import (
	"encoding/json"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

type Answers []db.AnswerItem

func (a Answers) String() string {
	val, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(val)
}

func NewAnswers(arr []db.AnswerItem) (*Answers, error) {
	ans := Answers(arr) // Convert slice to Answers type
	return &ans, nil    // Return pointer to Answers
}
