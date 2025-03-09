package question

import (
	"encoding/json"
	"go-rest-api-boilerplate/types"
)

type Answers []types.AnswerItem

func (a Answers) String() string {
	val, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(val)
}

func NewAnswers(arr []types.AnswerItem) (*Answers, error) {
	ans := Answers(arr) // Convert slice to Answers type
	return &ans, nil    // Return pointer to Answers
}
