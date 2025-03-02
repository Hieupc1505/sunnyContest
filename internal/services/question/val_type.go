package question

import "fmt"

type AnswerType string

type QuestionType string

func (q *AnswerType) String() string {
	return string(*q)
}

var validTypes = map[string]struct{}{
	"MULTI_CHOICE": {},
}

var AnswerTypes = map[string]struct{}{
	"TEXT":  {},
	"IMAGE": {},
}

func NewAnswerType(s string) (AnswerType, error) {
	_, exists := AnswerTypes[s]
	if !exists {
		return "", fmt.Errorf("invalid answer type: %s", s)
	}
	return AnswerType(s), nil
}

func (q QuestionType) String() string {
	return string(q)
}
func NewQuestionType(q string) (QuestionType, error) {
	_, exists := validTypes[q]
	if !exists {
		return "", fmt.Errorf("invalid question type: %s", q)
	}
	return QuestionType(q), nil
}
