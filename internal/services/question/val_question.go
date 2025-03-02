package question

import (
	"fmt"
)

const (
	MinQLen = 10
	MaxQLen = 1000
)

type Question string

func (q *Question) String() string {
	return string(*q)
}

func NewQuestion(s string) (Question, error) {
	if s == "" {
		return "", fmt.Errorf("Question cannot be empty")
	}
	if len(s) < MinQLen || len(s) > MaxQLen {
		return "", fmt.Errorf("Question length should be between %d and %d", MinQLen, MaxQLen)
	}
	return Question(s), nil
}
